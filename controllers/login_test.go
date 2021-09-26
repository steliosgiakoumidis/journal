package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

var createTokenResponse error

type FakeAuth struct {
}

func (f FakeAuth) CreateToken(w http.ResponseWriter, userId string) error {
	return createTokenResponse
}

func TestLoginFailsWhenValidateTokenReturnsError(t *testing.T) {

	createTokenResponse = errors.New("Sonething")
	handler := NewLoginHandler(FakeAuth{})
	req := httptest.NewRequest("GET", "/token", nil)
	w := httptest.NewRecorder()
	handler.GetToken(w, req)

	if w.Code != 500 {
		t.Errorf("Test")
	}
}

func TestLoginSucceedsWhenCreateTokenSuccesful(t *testing.T) {

	createTokenResponse = nil
	handler := NewLoginHandler(FakeAuth{})
	req := httptest.NewRequest("GET", "/token", nil)
	w := httptest.NewRecorder()
	handler.GetToken(w, req)

	if w.Code != 200 {
		t.Errorf("Test")
	}
}
