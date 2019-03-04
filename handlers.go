package main

import (
	"fmt"
	"net/http"
)

func HandleAuth(w http.ResponseWriter, r *http.Request, session *Session) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	if session.user != nil {
		fmt.Fprintln(w, "Already logged in")
	} else {
		user, err := Auth(login, password)
		if err != nil {
			fmt.Fprintln(w, "Auth error", err.Error())
		} else {
			session.user = user
			fmt.Fprintln(w, "Authorized", login, session.sid)
		}

	}
}

// HandleRegister handle registration api request
// request must contain post form:
// 	login, password, email, name
// Writes status json to response
func HandleRegister(w http.ResponseWriter, r *http.Request, session *Session) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	email := r.FormValue("email")
	name := r.FormValue("name")

	user, err := NewUser(login, password, email, name)
	if err == nil {
		session.user = user
		fmt.Fprintln(w, "User registered:", user.login, session.sid)
	} else {
		fmt.Fprintln(w, err.Error())
	}
}
