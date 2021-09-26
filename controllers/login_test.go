package controllers

import (
	"errors"
	"net/http"
	"testing"

	"net/http/httptest"
)

var createTokenResponse error

type FakeAuth struct {
}

func (f FakeAuth) CreateToken(w http.ResponseWriter, userId string) error {
	return createTokenResponse
}

func LoginFailsWhenValidateTokenReturnsError(t *testing.T) {

	createTokenResponse = errors.New("Sonething")
	handler := NewLoginHandler(FakeAuth{})
	req := httptest.NewRequest("GET", "/token", nil)
	w := httptest.NewRecorder()
	handler.GetToken(w, req)

	println(w.Code)

	if w.Code == 200 {
		t.Error("Test")
	}
}
