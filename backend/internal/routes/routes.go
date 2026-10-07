package routes

import (
	"net/http"
	"qare/backend/internal/handlers"
)

func New(a *handlers.API, origin string) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("POST /api/auth/register", a.Register)
	m.HandleFunc("POST /api/auth/login", a.Login)
	m.HandleFunc("GET /api/auth/me", a.Require(a.Me))
	m.HandleFunc("POST /api/auth/logout", a.Require(a.Logout))
	m.HandleFunc("GET /api/books", a.Books)
	m.HandleFunc("GET /api/books/{id}", a.Book)
	m.HandleFunc("PUT /api/books/{id}/review", a.Require(a.SaveReview))
	m.HandleFunc("DELETE /api/books/{id}/review", a.Require(a.DeleteReview))
	m.HandleFunc("GET /api/library", a.Require(a.Library))
	m.HandleFunc("PUT /api/library/{id}", a.Require(a.SaveLibrary))
	m.HandleFunc("DELETE /api/library/{id}", a.Require(a.DeleteLibrary))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		m.ServeHTTP(w, r)
	})
}
