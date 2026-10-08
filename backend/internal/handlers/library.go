package handlers

import (
	"errors"
	"net/http"
	"qare/backend/internal/repositories"
	"qare/backend/internal/services"
)

func (a *API) Library(w http.ResponseWriter, r *http.Request) {
	user, _ := a.current(r)
	books, err := repositories.Library(a.DB, user.ID)
	if err != nil {
		fail(w, 500, "تعذر تحميل المكتبة")
		return
	}
	send(w, 200, books)
}

func (a *API) SaveLibrary(w http.ResponseWriter, r *http.Request) {
	user, _ := a.current(r)
	var input struct {
		Status string `json:"status"`
	}
	if decode(r, &input) != nil {
		fail(w, 400, "تصنيف غير صالح")
		return
	}
	err := a.LibraryService.Save(id(r), user.ID, input.Status)
	switch {
	case errors.Is(err, services.ErrBookNotFound):
		fail(w, 404, "الكتاب غير موجود")
	case errors.Is(err, services.ErrInvalidStatus):
		fail(w, 400, "تصنيف غير صالح")
	case err != nil:
		fail(w, 500, "تعذر حفظ الكتاب")
	default:
		send(w, 200, map[string]string{"status": input.Status})
	}
}

func (a *API) DeleteLibrary(w http.ResponseWriter, r *http.Request) {
	user, _ := a.current(r)
	if err := a.LibraryService.Delete(id(r), user.ID); err != nil {
		fail(w, 500, "تعذر حذف الكتاب")
		return
	}
	w.WriteHeader(204)
}
