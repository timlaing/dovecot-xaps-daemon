package internal

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func newTestSocketHandler(t *testing.T) (*httpHandler, *httptest.Server) {
	t.Helper()
	apns, server := newAPNSMockClient(t, http.StatusOK)
	return &httpHandler{
		db:   newInMemoryDB(t),
		apns: apns,
	}, server
}

func performRequest(handler func(http.ResponseWriter, *http.Request, httprouter.Params), body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler(rec, req, httprouter.Params{})
	return rec
}

func TestHandleRegisterInvalidJSON(t *testing.T) {
	h, server := newTestSocketHandler(t)
	defer server.Close()

	rec := performRequest(h.handleRegister, []byte("{not json"))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleRegisterIncompleteParams(t *testing.T) {
	h, server := newTestSocketHandler(t)
	defer server.Close()

	body := []byte(`{"ApsAccountId":"acc","ApsDeviceToken":"tok","Username":"user"}`)
	rec := performRequest(h.handleRegister, body)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleRegisterWrongSubtopic(t *testing.T) {
	h, server := newTestSocketHandler(t)
	defer server.Close()

	body := []byte(`{"ApsAccountId":"acc","ApsDeviceToken":"tok","ApsSubtopic":"com.apple.wrong","Username":"user","Mailboxes":["Inbox"]}`)
	rec := performRequest(h.handleRegister, body)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleRegisterSuccess(t *testing.T) {
	h, server := newTestSocketHandler(t)
	defer server.Close()
	h.apns.Topic = "com.apple.mail.test"

	body := []byte(`{"ApsAccountId":"acc","ApsDeviceToken":"tok","ApsSubtopic":"com.apple.mobilemail","Username":"user","Mailboxes":["Inbox"]}`)
	rec := performRequest(h.handleRegister, body)
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != "com.apple.mail.test" {
		t.Errorf("body = %q, want %q", got, "com.apple.mail.test")
	}
}

func TestHandleNotifyInvalidJSON(t *testing.T) {
	h, server := newTestSocketHandler(t)
	defer server.Close()

	rec := performRequest(h.handleNotify, []byte("{not json"))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleNotifyIncompleteParams(t *testing.T) {
	h, server := newTestSocketHandler(t)
	defer server.Close()

	rec := performRequest(h.handleNotify, []byte(`{"Username":"user"}`))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleNotifyNonInbox(t *testing.T) {
	h, server := newTestSocketHandler(t)
	defer server.Close()

	rec := performRequest(h.handleNotify, []byte(`{"Username":"user","Mailbox":"Drafts","Events":["MessageNew"]}`))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestHandleNotifyNoRegistrationsUserExists(t *testing.T) {
	h, server := newTestSocketHandler(t)
	defer server.Close()

	rec := performRequest(h.handleNotify, []byte(`{"Username":"user","Mailbox":"INBOX","Events":["MessageNew"]}`))
	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestHandleNotifyNoRegistrationsUserMissing(t *testing.T) {
	h, server := newTestSocketHandler(t)
	defer server.Close()

	rec := performRequest(h.handleNotify, []byte(`{"Username":"missing","Mailbox":"INBOX","Events":["MessageNew"]}`))
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestHandleNotifyMessageNew(t *testing.T) {
	h, server := newTestSocketHandler(t)
	defer server.Close()

	db := newInMemoryDB(t)
	h.db = db
	// register an account whose mailbox matches "INBOX" (case-sensitive match)
	if err := db.AddRegistration("user", "acc-inbox", "tok", []string{"INBOX"}); err != nil {
		t.Fatal(err)
	}

	rec := performRequest(h.handleNotify, []byte(`{"Username":"user","Mailbox":"INBOX","Events":["MessageNew"]}`))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestHandleNotifyDelayed(t *testing.T) {
	h, server := newTestSocketHandler(t)
	defer server.Close()

	db := newInMemoryDB(t)
	h.db = db
	if err := db.AddRegistration("user", "acc-inbox", "tok", []string{"INBOX"}); err != nil {
		t.Fatal(err)
	}

	rec := performRequest(h.handleNotify, []byte(`{"Username":"user","Mailbox":"INBOX","Events":["MailboxCreate"]}`))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if len(h.apns.delayedApns) != 1 {
		t.Errorf("expected one delayed registration, got %d", len(h.apns.delayedApns))
	}
}

func TestRegisterCheckParams(t *testing.T) {
	tests := []struct {
		name  string
		reg   Register
		error bool
	}{
		{"all present", Register{ApsAccountId: "a", ApsDeviceToken: "b", Username: "c", Mailboxes: []string{"Inbox"}}, false},
		{"missing account id", Register{ApsDeviceToken: "b", Username: "c", Mailboxes: []string{"Inbox"}}, true},
		{"missing device token", Register{ApsAccountId: "a", Username: "c", Mailboxes: []string{"Inbox"}}, true},
		{"missing username", Register{ApsAccountId: "a", ApsDeviceToken: "b", Mailboxes: []string{"Inbox"}}, true},
		{"missing mailboxes", Register{ApsAccountId: "a", ApsDeviceToken: "b", Username: "c"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reg.checkParams(); got != tt.error {
				t.Errorf("checkParams() = %v, want %v", got, tt.error)
			}
		})
	}
}

func TestNotifyCheckParams(t *testing.T) {
	tests := []struct {
		name   string
		notify Notify
		error  bool
	}{
		{"all present", Notify{Username: "u", Mailbox: "INBOX", Events: []string{"MessageNew"}}, false},
		{"missing username", Notify{Mailbox: "INBOX", Events: []string{"MessageNew"}}, true},
		{"missing mailbox", Notify{Username: "u", Events: []string{"MessageNew"}}, true},
		{"missing events", Notify{Username: "u", Mailbox: "INBOX"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.notify.checkParams(); got != tt.error {
				t.Errorf("checkParams() = %v, want %v", got, tt.error)
			}
		})
	}
}
