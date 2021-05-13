package common

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

const (
	session_key = "session_key"
)

func GetSession(r *http.Request, key string) interface{} {
	session, _ := store.Get(r, session_key)
	return session.Values[key]
}

func SetSession(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {
	session, _ := store.Get(r, session_key)
	session.Values[key] = value
	return session.Save(r, w)
}
