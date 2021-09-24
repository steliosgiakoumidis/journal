package auth

import (
	"context"
	"net/http"
	"time"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userId string
		var exp time.Time
		var err error
		if userId, exp, err = ValidateToken(r); err != nil {
			http.Error(w, http.StatusText(401), 401)
		}

		//Refresh token if token is about to expire
		if exp.Sub(time.Now().UTC()) < time.Hour * 2 {
			CreateToken(w, userId)
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
