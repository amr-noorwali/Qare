package repositories

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenRejectsInvalidDSN(t *testing.T) {
	if _, err := Open("not a MySQL DSN"); err == nil {
		t.Fatal("expected invalid DSN error")
	}
}

func TestOpenRejectsInvalidCACertificate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(path, []byte("not a certificate"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MYSQL_CA_CERT_PATH", path)
	if _, err := Open("user:pass@tcp(localhost:3306)/qare"); err == nil {
		t.Fatal("expected invalid CA certificate error")
	}
}

func TestMySQLConnectionUsesUTF8MB4(t *testing.T) {
	dsn := os.Getenv("MYSQL_TEST_DSN")
	if dsn == "" {
		t.Skip("MYSQL_TEST_DSN is not set")
	}
	db, err := OpenForImport(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var charset string
	if err := db.QueryRow("SELECT @@character_set_connection").Scan(&charset); err != nil {
		t.Fatal(err)
	}
	if charset != "utf8mb4" {
		t.Fatalf("connection charset = %q, want utf8mb4", charset)
	}
}

// MYSQL_TEST_DSN must point to an isolated disposable MySQL database.
func TestMySQLSchemaAndSeed(t *testing.T) {
	dsn := os.Getenv("MYSQL_TEST_DSN")
	if dsn == "" {
		t.Skip("MYSQL_TEST_DSN is not set")
	}
	db, err := Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var books, version int
	if err := db.QueryRow("SELECT COUNT(*) FROM books").Scan(&books); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT MAX(version) FROM schema_migrations").Scan(&version); err != nil {
		t.Fatal(err)
	}
	if books < 6 || version != len(migrations) {
		t.Fatalf("schema/seed mismatch: books=%d migrations=%d", books, version)
	}
}
