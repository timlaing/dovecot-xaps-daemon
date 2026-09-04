package database

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDatabase_NewDatabaseMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "database.json")
	db, err := NewDatabase(path)
	if err != nil {
		t.Fatal("expected new database to be created, got error:", err)
	}
	if db.filename != path {
		t.Errorf("db.filename = %q, want %q", db.filename, path)
	}
	if len(db.Users) != 0 {
		t.Errorf("expected empty Users for new database, got %d", len(db.Users))
	}
}

func TestDatabase_NewDatabaseWriteError(t *testing.T) {
	// the parent directory does not exist and its parent cannot be created,
	// so the initial write() fails deterministically.
	dir := t.TempDir()
	_, err := NewDatabase(filepath.Join(dir, "missing-subdir", "database.json"))
	if err == nil {
		t.Fatal("expected error creating database in nonexistent directory")
	}
}

func TestDatabase_NewDatabaseReadError(t *testing.T) {
	// a directory is stat-able (so not IsNotExist) but not readable as a file
	dir := t.TempDir()

	if _, err := NewDatabase(dir); err == nil {
		t.Fatal("expected error reading a directory as a database file")
	}
}

func TestDatabase_NewDatabaseInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "database.json")
	if err := os.WriteFile(path, []byte("{ not valid json"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := NewDatabase(path); err == nil {
		t.Fatal("expected error unmarshalling invalid JSON")
	}
}

func TestDatabase_NewDatabaseEmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "database.json")
	if err := os.WriteFile(path, nil, 0644); err != nil {
		t.Fatal(err)
	}

	db, err := NewDatabase(path)
	if err != nil {
		t.Fatal("expected empty file to load, got error:", err)
	}
	if len(db.Users) != 0 {
		t.Errorf("expected empty Users from empty file, got %d", len(db.Users))
	}
}

func TestDatabase_AddRegistrationExistingAccount(t *testing.T) {
	path := filepath.Join(t.TempDir(), "database.json")
	db, err := NewDatabase(path)
	if err != nil {
		t.Fatal(err)
	}

	mailboxes := []string{"Inbox"}
	if err := db.AddRegistration("user@example.com", "account1", "token1", mailboxes); err != nil {
		t.Fatal(err)
	}
	// second registration for the same user + account id exercises the
	// "account already exists" branch
	if err := db.AddRegistration("user@example.com", "account1", "token2", mailboxes); err != nil {
		t.Fatal(err)
	}

	if got := db.Users["user@example.com"].Accounts["account1"].DeviceToken; got != "token2" {
		t.Errorf("DeviceToken = %q, want %q", got, "token2")
	}
}

func TestDatabase_DeleteIfExistRegistrationWriteError(t *testing.T) {
	// Construct a database whose backing file cannot be written so the
	// write() inside DeleteIfExistRegistration fails.
	db := &Database{
		filename: filepath.Join(t.TempDir(), "missing-subdir", "database.json"),
		Users: map[string]User{
			"user": {
				Accounts: map[string]Account{
					"account1": {DeviceToken: "token1", Mailboxes: []string{"Inbox"}},
				},
			},
		},
	}

	ok := db.DeleteIfExistRegistration(Registration{DeviceToken: "token1", AccountId: "account1"})
	if !ok {
		t.Fatal("expected registration to be deleted")
	}
	// after deletion the empty user is also cleaned up
	if _, exists := db.Users["user"]; exists {
		t.Error("expected empty user to be removed")
	}
}

func TestDatabase_UserExists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "database.json")
	db, err := NewDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	if db.UserExists("nobody@example.com") {
		t.Error("UserExists returned true for a user that does not exist")
	}

	if err := db.AddRegistration("someone@example.com", "account1", "token1", []string{"Inbox"}); err != nil {
		t.Fatal(err)
	}
	if !db.UserExists("someone@example.com") {
		t.Error("UserExists returned false for a user that exists")
	}
}

func TestDatabase_CleanupRegisteredRemovesStale(t *testing.T) {
	path := filepath.Join(t.TempDir(), "database.json")
	db := &Database{
		filename: path,
		Users: map[string]User{
			"user": {
				Accounts: map[string]Account{
					"fresh": {DeviceToken: "fresh", RegistrationTime: time.Now()},
					"old":   {DeviceToken: "old", RegistrationTime: time.Now().Add(-time.Hour * 24 * 31)},
				},
			},
		},
	}

	db.cleanupRegistered()

	if _, exists := db.Users["user"].Accounts["old"]; exists {
		t.Error("expected stale registration to be removed")
	}
	if _, exists := db.Users["user"].Accounts["fresh"]; !exists {
		t.Error("expected fresh registration to remain")
	}
	if len(db.Users["user"].Accounts) != 1 {
		t.Errorf("expected only the fresh account to remain, got %d", len(db.Users["user"].Accounts))
	}
}
