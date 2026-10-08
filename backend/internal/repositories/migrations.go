package repositories

import (
	"database/sql"
	"fmt"
)

// Append new migrations here; never edit an applied migration.
var migrations = []string{
	`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL UNIQUE, password_hash TEXT NOT NULL);
 CREATE TABLE IF NOT EXISTS sessions (token TEXT PRIMARY KEY, user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, expires_at TEXT NOT NULL);
 CREATE TABLE IF NOT EXISTS books (id INTEGER PRIMARY KEY, title TEXT NOT NULL, author TEXT NOT NULL, description TEXT NOT NULL, cover_url TEXT NOT NULL, genre TEXT NOT NULL);
 CREATE TABLE IF NOT EXISTS reviews (id INTEGER PRIMARY KEY, book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE, user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, rating INTEGER NOT NULL CHECK(rating BETWEEN 1 AND 5), body TEXT NOT NULL, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, UNIQUE(book_id,user_id));
 CREATE TABLE IF NOT EXISTS library (user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE, status TEXT NOT NULL CHECK(status IN ('want','reading','read')), PRIMARY KEY(user_id,book_id));
 CREATE INDEX IF NOT EXISTS idx_reviews_book_created ON reviews(book_id,created_at DESC);
 CREATE INDEX IF NOT EXISTS idx_library_book ON library(book_id);
 CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);`,
}

func migrate(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var version int
	if err := tx.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return err
	}
	if version > len(migrations) {
		return fmt.Errorf("database schema version %d is newer than supported version %d", version, len(migrations))
	}
	for i := version; i < len(migrations); i++ {
		if _, err := tx.Exec(migrations[i]); err != nil {
			return fmt.Errorf("migration %d: %w", i+1, err)
		}
		if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version=%d", i+1)); err != nil {
			return err
		}
	}
	return tx.Commit()
}
