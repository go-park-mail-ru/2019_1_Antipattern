package main

import "net/http"

func SessionMiddleware(next func(http.ResponseWriter, *http.Request, *Session)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sid")
		if err != nil {
			session := NewSession()
			cookie = &http.Cookie{
				Name:     "sid",
				Value:    session.sid,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}
		session, ok := sessions[cookie.Value]
		if !ok {
			session = *NewSession()
			cookie = &http.Cookie{
				Name:     "sid",
				Value:    session.sid,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}
		next(w, r, &session)
		session.Save()

	}
}
