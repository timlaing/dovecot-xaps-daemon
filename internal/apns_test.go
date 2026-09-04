package internal

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/sideshow/apns2"
	log "github.com/sirupsen/logrus"
	"github.com/timlaing/dovecot-xaps-daemon/internal/database"
)

// newSelfSignedCert builds an x509 certificate with the given subject
// names and returns it as tls.Certificate.
func newSelfSignedCert(t *testing.T, names []pkix.AttributeTypeAndValue) tls.Certificate {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		Subject:               pkix.Name{ExtraNames: names},
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatal(err)
	}
	return tlsCert
}

func TestTopicFromCertificate(t *testing.T) {
	t.Run("multiple certificates is an error", func(t *testing.T) {
		cert := newSelfSignedCert(t, []pkix.AttributeTypeAndValue{
			{Type: oidUid, Value: "uid1"},
			{Type: oidUid, Value: "uid2"},
		})
		// force two leaf certificates in the chain
		cert.Certificate = append(cert.Certificate, cert.Certificate[0])
		if _, err := topicFromCertificate(cert); err == nil {
			t.Fatal("expected error for multiple certificates")
		}
	})

	t.Run("empty subject names is an error", func(t *testing.T) {
		cert := newSelfSignedCert(t, nil)
		if _, err := topicFromCertificate(cert); err == nil {
			t.Fatal("expected error for empty subject names")
		}
	})

	t.Run("wrong subject name OID is an error", func(t *testing.T) {
		cert := newSelfSignedCert(t, []pkix.AttributeTypeAndValue{
			{Type: []int{2, 5, 4, 3}, Value: "common-name"},
		})
		if _, err := topicFromCertificate(cert); err == nil {
			t.Fatal("expected error for wrong subject name OID")
		}
	})

	t.Run("returns uid for a valid certificate", func(t *testing.T) {
		want := "com.apple.mail.test"
		cert := newSelfSignedCert(t, []pkix.AttributeTypeAndValue{
			{Type: oidUid, Value: want},
		})
		got, err := topicFromCertificate(cert)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("topic = %q, want %q", got, want)
		}
	})
}

func newInMemoryDB(t *testing.T) *database.Database {
	t.Helper()
	path := filepath.Join(t.TempDir(), "database.json")

	db, err := database.NewDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AddRegistration("user", "acc", "token", []string{"Inbox"}); err != nil {
		t.Fatal(err)
	}
	return db
}

func newTestApns(client *apns2.Client) *Apns {
	return &Apns{
		DelayTime:            1,
		CheckDelayedInterval: 1,
		client:               client,
		delayedApns:          make(map[database.Registration]time.Time),
	}
}

func newAPNSMockClient(t *testing.T, statusCode int) (*Apns, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))
	client := &apns2.Client{Host: server.URL, HTTPClient: http.DefaultClient}
	return newTestApns(client), server
}

func TestSendNotificationDelayed(t *testing.T) {
	apns, server := newAPNSMockClient(t, http.StatusOK)
	defer server.Close()

	apns.SendNotification(database.Registration{DeviceToken: "token", AccountId: "acc"}, true)

	if len(apns.delayedApns) != 1 {
		t.Fatalf("expected one delayed registration, got %d", len(apns.delayedApns))
	}
	if _, ok := apns.delayedApns[database.Registration{DeviceToken: "token", AccountId: "acc"}]; !ok {
		t.Error("registration not in delayed map")
	}
}

func TestSendNotificationStatusOK(t *testing.T) {
	apns, server := newAPNSMockClient(t, http.StatusOK)
	defer server.Close()
	apns.db = newInMemoryDB(t)

	apns.SendNotification(database.Registration{DeviceToken: "token", AccountId: "acc"}, false)
	if len(apns.delayedApns) != 0 {
		t.Errorf("delayed map should be empty, got %d", len(apns.delayedApns))
	}
}

func TestSendNotificationStatusUnregistered(t *testing.T) {
	apns, server := newAPNSMockClient(t, 410)
	defer server.Close()
	db := newInMemoryDB(t)
	apns.db = db

	apns.SendNotification(database.Registration{DeviceToken: "token", AccountId: "acc"}, false)

	if db.UserExists("user") {
		t.Error("expected registration to be deleted on 410")
	}
}

func TestSendNotificationStatusOther(t *testing.T) {
	apns, server := newAPNSMockClient(t, http.StatusServiceUnavailable)
	defer server.Close()
	apns.db = newInMemoryDB(t)

	apns.SendNotification(database.Registration{DeviceToken: "token", AccountId: "acc"}, false)
}

func TestCheckDelayed(t *testing.T) {
	apns, server := newAPNSMockClient(t, http.StatusOK)
	defer server.Close()
	apns.db = newInMemoryDB(t)

	apns.DelayTime = 10
	old := database.Registration{DeviceToken: "old-token", AccountId: "old"}
	recent := database.Registration{DeviceToken: "recent-token", AccountId: "recent"}
	apns.delayedApns[old] = time.Now().Add(-time.Hour)
	apns.delayedApns[recent] = time.Now()

	apns.checkDelayed()

	if _, ok := apns.delayedApns[old]; ok {
		t.Error("expected old registration to be sent and removed")
	}
	if _, ok := apns.delayedApns[recent]; !ok {
		t.Error("expected recent registration to remain")
	}
}

func TestSendNotificationDebugBranch(t *testing.T) {
	oldLevel := log.GetLevel()
	log.SetLevel(log.DebugLevel)
	defer log.SetLevel(oldLevel)

	apns, server := newAPNSMockClient(t, http.StatusOK)
	defer server.Close()
	apns.db = newInMemoryDB(t)

	apns.SendNotification(database.Registration{DeviceToken: "token", AccountId: "acc"}, false)
}

func TestCreateDelayedNotificationThread(t *testing.T) {
	apns, server := newAPNSMockClient(t, http.StatusOK)
	defer server.Close()
	apns.db = newInMemoryDB(t)
	apns.CheckDelayedInterval = 1

	apns.delayedApns[database.Registration{DeviceToken: "token", AccountId: "acc"}] = time.Now().Add(-time.Hour)

	apns.createDelayedNotificationThread()

	// let the one-second ticker fire so the goroutine body executes
	time.Sleep(1100 * time.Millisecond)

	// wait (with a timeout) for the delayed registration to be sent
	deadline := time.Now().Add(5 * time.Second)
	for {
		apns.mapMutex.Lock()
		n := len(apns.delayedApns)
		apns.mapMutex.Unlock()
		if n == 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for delayed registration to be sent")
		}
		time.Sleep(50 * time.Millisecond)
	}
}
