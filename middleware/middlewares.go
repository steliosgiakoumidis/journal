package middleware

import (
	"context"
	"journal/auth"
	"log"
	"net/http"
	"time"
)

func NewCustomMiddleware(auth auth.Auth) *CustomMiddleware {
	return &CustomMiddleware{
		auth: auth,
	}
}

type CustomMiddleware struct {
	auth auth.Auth
}

func (c CustomMiddleware) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userId string
		var exp time.Time
		var err error
		if userId, exp, err = c.auth.ValidateToken(r); err != nil {
			http.Error(w, http.StatusText(401), 401)
		}

		//Refresh token if token is about to expire
		if exp.Sub(time.Now().UTC()) < time.Hour*2 {
			c.auth.CreateToken(w, userId)
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c CustomMiddleware) TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeNow := time.Now()
		log.Println("Request started at:" + timeNow.String())
		next.ServeHTTP(w, r)

		log.Println("Request completed at: " + time.Now().String())

		timeAtTheEnd := time.Now()
		log.Println("Request lasted: " + timeAtTheEnd.Sub(timeNow).String())
	})
}
