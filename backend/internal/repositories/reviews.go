package repositories

import "database/sql"

type Review struct {
	ID        int64  `json:"id"`
	BookID    int64  `json:"book_id"`
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	Rating    int    `json:"rating"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

func Reviews(db *sql.DB, bookID int64) ([]Review, error) {
	rows, err := db.Query(`SELECT r.id,r.book_id,r.user_id,u.name,r.rating,r.body,r.created_at FROM reviews r JOIN users u ON u.id=r.user_id WHERE r.book_id=? ORDER BY r.created_at DESC,r.id DESC`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reviews := []Review{}
	for rows.Next() {
		var review Review
		if err := rows.Scan(&review.ID, &review.BookID, &review.UserID, &review.UserName, &review.Rating, &review.Body, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, rows.Err()
}

func SaveReview(db *sql.DB, bookID, userID int64, rating int, body string) error {
	_, err := db.Exec(`INSERT INTO reviews(book_id,user_id,rating,body) VALUES(?,?,?,?)
 ON DUPLICATE KEY UPDATE rating=VALUES(rating),body=VALUES(body),created_at=CURRENT_TIMESTAMP`, bookID, userID, rating, body)
	return err
}

func DeleteReview(db *sql.DB, bookID, userID int64) (bool, error) {
	result, err := db.Exec("DELETE FROM reviews WHERE book_id=? AND user_id=?", bookID, userID)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	return affected > 0, err
}
