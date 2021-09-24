package controllers

import (
	"journal/auth"
	"net/http"
)

func GetToken(w http.ResponseWriter, r *http.Request) {
	if err := auth.CreateToken(w, "123"); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(200)
}
