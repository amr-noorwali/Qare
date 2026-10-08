// Command importsqlite copies an existing Qare SQLite database into an empty
// MySQL database. Stop the old API before running it so the source cannot change.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"qare/backend/internal/repositories"
	"strings"

	_ "modernc.org/sqlite"
)

type table struct {
	name    string
	columns string
}

var tables = []table{
	{"users", "id,name,email,password_hash"},
	{"books", "id,title,author,description,cover_url,genre"},
	{"sessions", "token,user_id,expires_at"},
	{"reviews", "id,book_id,user_id,rating,body,created_at"},
	{"library", "user_id,book_id,status"},
}

func main() {
	if len(os.Args) != 2 || os.Getenv("MYSQL_DSN") == "" {
		log.Fatal("usage: MYSQL_DSN=<dsn> go run ./cmd/importsqlite <path-to-qare.db>")
	}
	if _, err := os.Stat(os.Args[1]); err != nil {
		log.Fatal(err)
	}
	source, err := sql.Open("sqlite", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()
	target, err := repositories.OpenForImport(os.Getenv("MYSQL_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer target.Close()
	if err := copyDatabase(source, target); err != nil {
		log.Fatal(err)
	}
	log.Println("SQLite data copied to MySQL successfully")
}

func copyDatabase(source, target *sql.DB) error {
	tx, err := target.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, t := range tables {
		var count int
		if err := tx.QueryRow("SELECT COUNT(*) FROM " + t.name).Scan(&count); err != nil {
			return err
		}
		if count != 0 {
			return fmt.Errorf("MySQL %s table is not empty; import requires a new database", t.name)
		}
	}
	for _, t := range tables {
		rows, err := source.Query("SELECT " + t.columns + " FROM " + t.name)
		if err != nil {
			return fmt.Errorf("read SQLite %s: %w", t.name, err)
		}
		columnCount := len(strings.Split(t.columns, ","))
		query := "INSERT INTO " + t.name + "(" + t.columns + ") VALUES(" + strings.TrimSuffix(strings.Repeat("?,", columnCount), ",") + ")"
		for rows.Next() {
			values := make([]any, columnCount)
			destinations := make([]any, columnCount)
			for i := range values {
				destinations[i] = &values[i]
			}
			if err := rows.Scan(destinations...); err != nil {
				rows.Close()
				return err
			}
			if _, err := tx.Exec(query, values...); err != nil {
				rows.Close()
				return fmt.Errorf("write MySQL %s: %w", t.name, err)
			}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return err
		}
		rows.Close()
	}
	return tx.Commit()
}
