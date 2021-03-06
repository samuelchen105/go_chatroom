package common

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
)

type CtxKey string

var (
	secureC *securecookie.SecureCookie
)

const (
	hashKey    = "D2D02EA74DE2C9FAB1D802DB969C18D4"
	blockKey   = "0F9539500E9DB6826453864CC0BE85BA"
	cookieName = "chatuser"
)

type UserCookie struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func InitCookie() {
	secureC = securecookie.New([]byte(hashKey), []byte(blockKey))
}

func SetCookie(w http.ResponseWriter, data *UserCookie) error {
	cval, err := secureC.Encode(cookieName, data)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  cval,
		MaxAge: 0,
		Path:   "/",
	}

	http.SetCookie(w, cookie)
	return nil
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		cval := &UserCookie{}
		if err := secureC.Decode(cookieName, cookie.Value, cval); err != nil {
			log.Println("decode secure cookie: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newRequest := r.WithContext(context.WithValue(r.Context(), CtxKey(cookieName), cval))

		next.ServeHTTP(w, newRequest)
	})
}

func ReadCookie(w http.ResponseWriter, r *http.Request) (*UserCookie, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, errors.New("user do not login")
	}

	cval := &UserCookie{}
	if err := secureC.Decode(cookieName, cookie.Value, cval); err != nil {
		return nil, err
	}

	return cval, nil
}
