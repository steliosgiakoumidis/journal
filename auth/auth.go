package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"time"
)

type Auth struct {
}

func (a Auth) ValidateToken(r *http.Request) (string, time.Time, error) {
	var claims jwt.MapClaims
	var ok bool
	token, err := verifyToken(r)
	if err != nil {
		return "", time.Time{}, err
	}

	if claims, ok = token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		err = errors.New("token invalid or malformed")
		log.Println(err)
		return "", time.Time{}, err
	}

	userId := claims["user_id"].(string)
	expD, err := time.Parse(time.RFC3339, claims["exp"].(string))
	if err != nil {
		err = errors.New("token invalid or malformed")
		log.Println(err)
		return "", time.Time{}, err
	}

	return userId, expD, nil
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			test := token
			println(test)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		os.Setenv("ACCESS_SECRET", "test") //this should be in an env file
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func extractToken(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		log.Println("Cookie failed")
	}

	//normally Authorization the_token_xxx
	return cookie.Value
}

func (a Auth) CreateToken(w http.ResponseWriter, userid string) error {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "test") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Hour * 6).UTC()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return nil
}
