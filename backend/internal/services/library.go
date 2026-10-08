package services

import (
	"database/sql"
	"errors"
	"qare/backend/internal/repositories"
)

var ErrInvalidStatus = errors.New("invalid library status")

type Library struct{ DB *sql.DB }

func (s Library) Save(bookID, userID int64, status string) error {
	if err := requireBook(s.DB, bookID); err != nil {
		return err
	}
	if status != "want" && status != "reading" && status != "read" {
		return ErrInvalidStatus
	}
	return repositories.SaveLibrary(s.DB, userID, bookID, status)
}

func (s Library) Delete(bookID, userID int64) error {
	return repositories.DeleteLibrary(s.DB, userID, bookID)
}
