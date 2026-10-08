package services

import (
	"database/sql"
	"errors"
	"qare/backend/internal/repositories"
	"strings"
)

var ErrInvalidReview = errors.New("invalid review")
var ErrReviewNotFound = errors.New("review not found")

type Reviews struct{ DB *sql.DB }

func (s Reviews) Save(bookID, userID int64, rating int, body string) error {
	if err := requireBook(s.DB, bookID); err != nil {
		return err
	}
	body = strings.TrimSpace(body)
	if rating < 1 || rating > 5 || len(body) < 3 {
		return ErrInvalidReview
	}
	return repositories.SaveReview(s.DB, bookID, userID, rating, body)
}

func (s Reviews) Delete(bookID, userID int64) error {
	deleted, err := repositories.DeleteReview(s.DB, bookID, userID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrReviewNotFound
	}
	return nil
}
