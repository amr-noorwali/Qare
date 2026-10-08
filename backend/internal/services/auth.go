package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"qare/backend/internal/repositories"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	id, err := repositories.CreateUser(a.DB, name, email, string(hashed))
	if err != nil {
		return User{}, "", errors.New("البريد الإلكتروني مستخدم بالفعل")
	}
	user := User{id, name, email}
	token, err := a.session(id)
	return user, token, err
}

func (a Auth) Login(email, password string) (User, string, error) {
	record, err := repositories.UserByEmail(a.DB, strings.ToLower(strings.TrimSpace(email)))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(record.PasswordHash), []byte(password)) != nil {
		return User{}, "", ErrInvalid
	}
	token, err := a.session(record.ID)
	return User{record.ID, record.Name, record.Email}, token, err
}

func (a Auth) session(userID int64) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)
	err := repositories.CreateSession(a.DB, token, userID, time.Now().Add(30*24*time.Hour).UTC().Format(time.RFC3339))
	return token, err
}

func (a Auth) User(token string) (User, error) {
	record, err := repositories.UserBySession(a.DB, token, time.Now().UTC().Format(time.RFC3339))
	return User{record.ID, record.Name, record.Email}, err
}

func (a Auth) Logout(token string) { repositories.DeleteSession(a.DB, token) }
