package repositories

import (
	"path/filepath"
	"testing"
)

func TestOpenUpgradesExistingDatabaseWithoutChangingData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "books.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("INSERT INTO users(name,email,password_hash) VALUES('مستخدم','user@example.com','hash')"); err != nil {
		t.Fatal(err)
	}
	// Existing MVP databases have the tables but no schema version.
	if _, err := db.Exec("PRAGMA user_version=0"); err != nil {
		t.Fatal(err)
	}
	db.Close()

	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var version, books, users int
	if err := db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM books").Scan(&books); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email='user@example.com'").Scan(&users); err != nil {
		t.Fatal(err)
	}
	if version != len(migrations) || books != 6 || users != 1 {
		t.Fatalf("reopen changed data: version=%d books=%d users=%d", version, books, users)
	}
}
