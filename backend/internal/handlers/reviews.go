package handlers

import (
	"errors"
	"net/http"
	"qare/backend/internal/services"
)

func (a *API) SaveReview(w http.ResponseWriter, r *http.Request) {
	user, _ := a.current(r)
	var input struct {
		Rating int    `json:"rating"`
		Body   string `json:"body"`
	}
	if decode(r, &input) != nil {
		fail(w, 400, "أدخل تقييمًا من 1 إلى 5 ومراجعة نصية")
		return
	}
	err := a.ReviewsService.Save(id(r), user.ID, input.Rating, input.Body)
	switch {
	case errors.Is(err, services.ErrBookNotFound):
		fail(w, 404, "الكتاب غير موجود")
	case errors.Is(err, services.ErrInvalidReview):
		fail(w, 400, "أدخل تقييمًا من 1 إلى 5 ومراجعة نصية")
	case err != nil:
		fail(w, 500, "تعذر حفظ المراجعة")
	default:
		send(w, 200, map[string]string{"message": "تم حفظ المراجعة"})
	}
}

func (a *API) DeleteReview(w http.ResponseWriter, r *http.Request) {
	user, _ := a.current(r)
	err := a.ReviewsService.Delete(id(r), user.ID)
	switch {
	case errors.Is(err, services.ErrReviewNotFound):
		fail(w, 404, "المراجعة غير موجودة")
	case err != nil:
		fail(w, 500, "تعذر حذف المراجعة")
	default:
		w.WriteHeader(204)
	}
}
