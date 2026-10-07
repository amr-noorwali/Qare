package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"qare/backend/internal/repositories"
	"qare/backend/internal/services"
	"strconv"
	"strings"
)

type API struct {
	DB   *sql.DB
	Auth services.Auth
}

func New(db *sql.DB) *API { return &API{DB: db, Auth: services.Auth{DB: db}} }
func send(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func fail(w http.ResponseWriter, status int, msg string) {
	send(w, status, map[string]string{"error": msg})
}
func decode(r *http.Request, v any) error {
	return json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(v)
}
func (a *API) current(r *http.Request) (services.User, error) {
	t := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if t == "" || t == r.Header.Get("Authorization") {
		return services.User{}, errors.New("unauthorized")
	}
	return a.Auth.User(t)
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
func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	var x struct{ Name, Email, Password string }
	if decode(r, &x) != nil {
		fail(w, 400, "طلب غير صالح")
		return
	}
	u, t, e := a.Auth.Register(x.Name, x.Email, x.Password)
	if e != nil {
		fail(w, 400, e.Error())
		return
	}
	send(w, 201, map[string]any{"user": u, "token": t})
}
func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var x struct{ Email, Password string }
	if decode(r, &x) != nil {
		fail(w, 400, "طلب غير صالح")
		return
	}
	u, t, e := a.Auth.Login(x.Email, x.Password)
	if e != nil {
		fail(w, 401, e.Error())
		return
	}
	send(w, 200, map[string]any{"user": u, "token": t})
}
func (a *API) Me(w http.ResponseWriter, r *http.Request) { u, _ := a.current(r); send(w, 200, u) }
func (a *API) Logout(w http.ResponseWriter, r *http.Request) {
	a.Auth.Logout(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	w.WriteHeader(204)
}
func id(r *http.Request) int64 { n, _ := strconv.ParseInt(r.PathValue("id"), 10, 64); return n }
func bookExists(a *API, w http.ResponseWriter, bookID int64) bool {
	if bookID < 1 {
		fail(w, 404, "الكتاب غير موجود")
		return false
	}
	if _, e := repositories.OneBook(a.DB, bookID); e != nil {
		fail(w, 404, "الكتاب غير موجود")
		return false
	}
	return true
}
func (a *API) Books(w http.ResponseWriter, r *http.Request) {
	b, e := repositories.Books(a.DB, r.URL.Query().Get("q"))
	if e != nil {
		fail(w, 500, "تعذر تحميل الكتب")
		return
	}
	send(w, 200, b)
}
func (a *API) Book(w http.ResponseWriter, r *http.Request) {
	b, e := repositories.OneBook(a.DB, id(r))
	if e != nil {
		fail(w, 404, "الكتاب غير موجود")
		return
	}
	reviews, e := repositories.Reviews(a.DB, b.ID)
	if e != nil {
		fail(w, 500, "تعذر تحميل المراجعات")
		return
	}
	send(w, 200, map[string]any{"book": b, "reviews": reviews})
}
func (a *API) SaveReview(w http.ResponseWriter, r *http.Request) {
	u, _ := a.current(r)
	bid := id(r)
	if !bookExists(a, w, bid) {
		return
	}
	var x struct {
		Rating int    `json:"rating"`
		Body   string `json:"body"`
	}
	if decode(r, &x) != nil || x.Rating < 1 || x.Rating > 5 || len(strings.TrimSpace(x.Body)) < 3 {
		fail(w, 400, "أدخل تقييمًا من 1 إلى 5 ومراجعة نصية")
		return
	}
	_, e := a.DB.Exec(`INSERT INTO reviews(book_id,user_id,rating,body) VALUES(?,?,?,?) ON CONFLICT(book_id,user_id) DO UPDATE SET rating=excluded.rating,body=excluded.body,created_at=CURRENT_TIMESTAMP`, bid, u.ID, x.Rating, strings.TrimSpace(x.Body))
	if e != nil {
		fail(w, 500, "تعذر حفظ المراجعة")
		return
	}
	send(w, 200, map[string]string{"message": "تم حفظ المراجعة"})
}
func (a *API) DeleteReview(w http.ResponseWriter, r *http.Request) {
	u, _ := a.current(r)
	res, e := a.DB.Exec("DELETE FROM reviews WHERE book_id=? AND user_id=?", id(r), u.ID)
	if e != nil {
		fail(w, 500, "تعذر حذف المراجعة")
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		fail(w, 404, "المراجعة غير موجودة")
		return
	}
	w.WriteHeader(204)
}
func (a *API) Library(w http.ResponseWriter, r *http.Request) {
	u, _ := a.current(r)
	b, e := repositories.Library(a.DB, u.ID)
	if e != nil {
		fail(w, 500, "تعذر تحميل المكتبة")
		return
	}
	send(w, 200, b)
}
func (a *API) SaveLibrary(w http.ResponseWriter, r *http.Request) {
	u, _ := a.current(r)
	bid := id(r)
	if !bookExists(a, w, bid) {
		return
	}
	var x struct {
		Status string `json:"status"`
	}
	if decode(r, &x) != nil || (x.Status != "want" && x.Status != "reading" && x.Status != "read") {
		fail(w, 400, "تصنيف غير صالح")
		return
	}
	_, e := a.DB.Exec(`INSERT INTO library(user_id,book_id,status) VALUES(?,?,?) ON CONFLICT(user_id,book_id) DO UPDATE SET status=excluded.status`, u.ID, bid, x.Status)
	if e != nil {
		fail(w, 500, "تعذر حفظ الكتاب")
		return
	}
	send(w, 200, map[string]string{"status": x.Status})
}
func (a *API) DeleteLibrary(w http.ResponseWriter, r *http.Request) {
	u, _ := a.current(r)
	a.DB.Exec("DELETE FROM library WHERE user_id=? AND book_id=?", u.ID, id(r))
	w.WriteHeader(204)
}
