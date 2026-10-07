package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type Auth struct{ DB *sql.DB }
type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var ErrInvalid = errors.New("بيانات الدخول غير صحيحة")

func (a Auth) Register(name, email, password string) (User, string, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	if len(name) < 2 || !strings.Contains(email, "@") || len(password) < 8 {
		return User{}, "", errors.New("تحقق من الاسم والبريد وكلمة المرور (8 أحرف على الأقل)")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, "", err
	}
	r, err := a.DB.Exec("INSERT INTO users(name,email,password_hash) VALUES(?,?,?)", name, email, string(hashed))
	if err != nil {
		return User{}, "", errors.New("البريد الإلكتروني مستخدم بالفعل")
	}
	id, _ := r.LastInsertId()
	u := User{id, name, email}
	t, err := a.session(id)
	return u, t, err
}
func (a Auth) Login(email, password string) (User, string, error) {
	var u User
	var h string
	err := a.DB.QueryRow("SELECT id,name,email,password_hash FROM users WHERE email=?", strings.ToLower(strings.TrimSpace(email))).Scan(&u.ID, &u.Name, &u.Email, &h)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(h), []byte(password)) != nil {
		return User{}, "", ErrInvalid
	}
	t, err := a.session(u.ID)
	return u, t, err
}
func (a Auth) session(id int64) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	t := hex.EncodeToString(b)
	_, err := a.DB.Exec("INSERT INTO sessions(token,user_id,expires_at) VALUES(?,?,?)", t, id, time.Now().Add(30*24*time.Hour).UTC().Format(time.RFC3339))
	return t, err
}
func (a Auth) User(token string) (User, error) {
	var u User
	err := a.DB.QueryRow("SELECT u.id,u.name,u.email FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token=? AND s.expires_at>?", token, time.Now().UTC().Format(time.RFC3339)).Scan(&u.ID, &u.Name, &u.Email)
	return u, err
}
func (a Auth) Logout(token string) { a.DB.Exec("DELETE FROM sessions WHERE token=?", token) }
