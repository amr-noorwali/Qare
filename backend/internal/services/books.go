package services

import (
	"database/sql"
	"errors"
	"qare/backend/internal/repositories"
)

var ErrBookNotFound = errors.New("book not found")

func requireBook(db *sql.DB, bookID int64) error {
	if bookID < 1 {
		return ErrBookNotFound
	}
	exists, err := repositories.BookExists(db, bookID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrBookNotFound
	}
	return nil
}
