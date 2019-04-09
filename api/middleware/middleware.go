package middleware

import (
	"net/http"

	"../models"
	_ "github.com/dgrijalva/jwt-go"
)

func SessionMiddleware(next func(http.ResponseWriter, *http.Request, *models.Session), authRequiered bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sid")
		if err != nil {
			session := models.NewSession()
			cookie = &http.Cookie{
				Name:     "sid",
				Value:    session.Sid,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}
		session, ok := models.Sessions[cookie.Value]
		if !ok {
			session = *models.NewSession()
			cookie = &http.Cookie{
				Name:     "sid",
				Value:    session.Sid,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}
		if authRequiered && session.User == nil {
			w.WriteHeader(http.StatusForbidden)
			session.Save()
			return
		}
		next(w, r, &session)
		session.Save()

	}
}
