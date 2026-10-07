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
type Review struct {
	ID        int64  `json:"id"`
	BookID    int64  `json:"book_id"`
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	Rating    int    `json:"rating"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}
type LibraryBook struct {
	Book
	Status string `json:"status"`
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
func Reviews(db *sql.DB, id int64) ([]Review, error) {
	rows, err := db.Query(`SELECT r.id,r.book_id,r.user_id,u.name,r.rating,r.body,r.created_at FROM reviews r JOIN users u ON u.id=r.user_id WHERE r.book_id=? ORDER BY r.created_at DESC,r.id DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Review{}
	for rows.Next() {
		var r Review
		if e := rows.Scan(&r.ID, &r.BookID, &r.UserID, &r.UserName, &r.Rating, &r.Body, &r.CreatedAt); e != nil {
			return nil, e
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
func Library(db *sql.DB, userID int64) ([]LibraryBook, error) {
	rows, err := db.Query(`SELECT b.id,b.title,b.author,b.description,b.cover_url,b.genre,COALESCE(ROUND(AVG(r.rating),1),0),COUNT(r.id),l.status FROM library l JOIN books b ON b.id=l.book_id LEFT JOIN reviews r ON r.book_id=b.id WHERE l.user_id=? GROUP BY b.id,l.status ORDER BY b.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []LibraryBook{}
	for rows.Next() {
		var b Book
		var status string
		if e := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.CoverURL, &b.Genre, &b.Average, &b.ReviewCount, &status); e != nil {
			return nil, e
		}
		out = append(out, LibraryBook{b, status})
	}
	return out, rows.Err()
}
