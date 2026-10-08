package handlers

import (
	"net/http"
	"strings"
)

func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	var input struct{ Name, Email, Password string }
	if decode(r, &input) != nil {
		fail(w, 400, "طلب غير صالح")
		return
	}
	user, token, err := a.Auth.Register(input.Name, input.Email, input.Password)
	if err != nil {
		fail(w, 400, err.Error())
		return
	}
	send(w, 201, map[string]any{"user": user, "token": token})
}

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var input struct{ Email, Password string }
	if decode(r, &input) != nil {
		fail(w, 400, "طلب غير صالح")
		return
	}
	user, token, err := a.Auth.Login(input.Email, input.Password)
	if err != nil {
		fail(w, 401, err.Error())
		return
	}
	send(w, 200, map[string]any{"user": user, "token": token})
}

func (a *API) Me(w http.ResponseWriter, r *http.Request) {
	user, _ := a.current(r)
	send(w, 200, user)
}

func (a *API) Logout(w http.ResponseWriter, r *http.Request) {
	a.Auth.Logout(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	w.WriteHeader(204)
}
