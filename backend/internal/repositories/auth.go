package repositories

import "database/sql"

type UserRecord struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string
}

func CreateUser(db *sql.DB, name, email, passwordHash string) (int64, error) {
	result, err := db.Exec("INSERT INTO users(name,email,password_hash) VALUES(?,?,?)", name, email, passwordHash)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UserByEmail(db *sql.DB, email string) (UserRecord, error) {
	var user UserRecord
	err := db.QueryRow("SELECT id,name,email,password_hash FROM users WHERE email=?", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	return user, err
}

func UserBySession(db *sql.DB, token, now string) (UserRecord, error) {
	var user UserRecord
	err := db.QueryRow("SELECT u.id,u.name,u.email FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token=? AND s.expires_at>?", token, now).
		Scan(&user.ID, &user.Name, &user.Email)
	return user, err
}

func CreateSession(db *sql.DB, token string, userID int64, expiresAt string) error {
	_, err := db.Exec("INSERT INTO sessions(token,user_id,expires_at) VALUES(?,?,?)", token, userID, expiresAt)
	return err
}

func DeleteSession(db *sql.DB, token string) {
	db.Exec("DELETE FROM sessions WHERE token=?", token)
}
