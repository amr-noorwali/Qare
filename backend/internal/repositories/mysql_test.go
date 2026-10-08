package repositories

import (
	"os"
	"testing"
)

func TestOpenRejectsInvalidDSN(t *testing.T) {
	if _, err := Open("not a MySQL DSN"); err == nil {
		t.Fatal("expected invalid DSN error")
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
