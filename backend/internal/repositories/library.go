package repositories

import "database/sql"

type LibraryBook struct {
	Book
	Status string `json:"status"`
}

func Library(db *sql.DB, userID int64) ([]LibraryBook, error) {
	rows, err := db.Query(`SELECT b.id,b.title,b.author,b.description,b.cover_url,b.genre,COALESCE(ROUND(AVG(r.rating),1),0),COUNT(r.id),l.status FROM library l JOIN books b ON b.id=l.book_id LEFT JOIN reviews r ON r.book_id=b.id WHERE l.user_id=? GROUP BY b.id,l.status ORDER BY b.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	books := []LibraryBook{}
	for rows.Next() {
		var book Book
		var status string
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.CoverURL, &book.Genre, &book.Average, &book.ReviewCount, &status); err != nil {
			return nil, err
		}
		books = append(books, LibraryBook{book, status})
	}
	return books, rows.Err()
}

func SaveLibrary(db *sql.DB, userID, bookID int64, status string) error {
	_, err := db.Exec(`INSERT INTO library(user_id,book_id,status) VALUES(?,?,?)
 ON CONFLICT(user_id,book_id) DO UPDATE SET status=excluded.status`, userID, bookID, status)
	return err
}

func DeleteLibrary(db *sql.DB, userID, bookID int64) error {
	_, err := db.Exec("DELETE FROM library WHERE user_id=? AND book_id=?", userID, bookID)
	return err
}
