package controllers

import (
	"net/http"
)

type Auth interface {
	CreateToken(w http.ResponseWriter, userid string) error
}

type LoginHandler struct {
	auth Auth
}

func NewLoginHandler(authentication Auth) *LoginHandler {
	return &LoginHandler{
		auth: authentication,
	}
}

func (l *LoginHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	if err := l.auth.CreateToken(w, "123"); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(200)
}
