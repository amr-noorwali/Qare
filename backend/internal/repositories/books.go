package repositories

import (
	"database/sql"
	"strings"
)

type Book struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	CoverURL    string  `json:"cover_url"`
	Genre       string  `json:"genre"`
	Average     float64 `json:"average_rating"`
	ReviewCount int     `json:"review_count"`
}

const bookSelect = `SELECT b.id,b.title,b.author,b.description,b.cover_url,b.genre,COALESCE(ROUND(AVG(r.rating),1),0),COUNT(r.id) FROM books b LEFT JOIN reviews r ON r.book_id=b.id`
const bookGroup = ` GROUP BY b.id`

func scanBook(rows *sql.Rows) (Book, error) {
	var b Book
	err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.CoverURL, &b.Genre, &b.Average, &b.ReviewCount)
	return b, err
}
func Books(db *sql.DB, q string) ([]Book, error) {
	rows, err := db.Query(bookSelect+` WHERE b.title LIKE ? OR b.author LIKE ?`+bookGroup+` ORDER BY b.id`, "%"+strings.TrimSpace(q)+"%", "%"+strings.TrimSpace(q)+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Book{}
	for rows.Next() {
		b, e := scanBook(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, b)
	}
	return out, rows.Err()
}
func OneBook(db *sql.DB, id int64) (Book, error) {
	rows, err := db.Query(bookSelect+` WHERE b.id=?`+bookGroup, id)
	if err != nil {
		return Book{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return Book{}, sql.ErrNoRows
	}
	return scanBook(rows)
}

func BookExists(db *sql.DB, id int64) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM books WHERE id=?)", id).Scan(&exists)
	return exists, err
}
