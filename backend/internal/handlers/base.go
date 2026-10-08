package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"qare/backend/internal/services"
	"strconv"
	"strings"
)

type API struct {
	DB             *sql.DB
	Auth           services.Auth
	ReviewsService services.Reviews
	LibraryService services.Library
}

func New(db *sql.DB) *API {
	return &API{
		DB:             db,
		Auth:           services.Auth{DB: db},
		ReviewsService: services.Reviews{DB: db},
		LibraryService: services.Library{DB: db},
	}
}

func send(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func fail(w http.ResponseWriter, status int, message string) {
	send(w, status, map[string]string{"error": message})
}

func decode(r *http.Request, value any) error {
	return json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(value)
}

func (a *API) current(r *http.Request) (services.User, error) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token == "" || token == r.Header.Get("Authorization") {
		return services.User{}, errors.New("unauthorized")
	}
	return a.Auth.User(token)
}

func (a *API) Require(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := a.current(r); err != nil {
			fail(w, 401, "سجّل الدخول أولًا")
			return
		}
		next(w, r)
	}
}

func id(r *http.Request) int64 {
	number, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return number
}
