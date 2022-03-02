package gohttpservice

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

const SessionIdName = "sessionid"

func ises(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{}
	c.Name = SessionIdName
	c.Value = sessionId()
	c.MaxAge = sessionDur

	http.SetCookie(w, &c)
	r.AddCookie(&c)

}

func hses(w http.ResponseWriter, r *http.Request) {
	c, e := r.Cookie(SessionIdName)

	if e != nil || c == nil {
		ises(w, r)
	}

}

func sessionMan(ha http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			hses(w, r)
			ha.ServeHTTP(w, r)
			return
		})

}
