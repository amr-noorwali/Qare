package repositories

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err = db.Exec("PRAGMA foreign_keys=ON"); err == nil {
		err = migrate(db)
	}
	if err == nil {
		err = seed(db)
	}
	if err == nil {
		_, err = db.Exec("UPDATE users SET name=? WHERE email=? AND password_hash=?", "قارئ من منصة قارئ", "demo@qare.local", "seed-only")
	}
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
