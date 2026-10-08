package handlers

import (
	"net/http"
	"qare/backend/internal/repositories"
)

func (a *API) Books(w http.ResponseWriter, r *http.Request) {
	books, err := repositories.Books(a.DB, r.URL.Query().Get("q"))
	if err != nil {
		fail(w, 500, "تعذر تحميل الكتب")
		return
	}
	send(w, 200, books)
}

func (a *API) Book(w http.ResponseWriter, r *http.Request) {
	book, err := repositories.OneBook(a.DB, id(r))
	if err != nil {
		fail(w, 404, "الكتاب غير موجود")
		return
	}
	reviews, err := repositories.Reviews(a.DB, book.ID)
	if err != nil {
		fail(w, 500, "تعذر تحميل المراجعات")
		return
	}
	send(w, 200, map[string]any{"book": book, "reviews": reviews})
}
