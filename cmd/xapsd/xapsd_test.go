package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestHashPasswordHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// hashPassword terminates the process via os.Exit(0)
	hashPassword()
}

func TestHashPassword(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestHashPasswordHelperProcess")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	cmd.Stdin = strings.NewReader("correct horse battery staple\n")

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	sum := sha256.Sum256([]byte("correct horse battery staple"))
	wantHash := hex.EncodeToString(sum[:])
	if !strings.Contains(out.String(), wantHash) {
		t.Errorf("output %q does not contain hash %q", out.String(), wantHash)
	}
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Fatal("Version must not be empty")
	}
}
