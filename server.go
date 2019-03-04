package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type User struct {
	uuid          uint32
	login         string
	password_hash string
	email         string
	name          string
}

type Session struct {
	sid  string
	user *User
}

func (session *Session) Save() error {
	sessions[session.sid] = *session
	return nil
}
func (user *User) Save() error {
	// TODO: Save to db logic
	return nil
}

var users map[string]User
var sessions map[string]Session

func NewSession() *Session {
	id := uuid.New().String()
	session := Session{
		sid:  id,
		user: nil,
	}

	sessions[id] = session
	return &session
}

func NewUser(login string, password string, email string, name string) (*User, error) {
	if _, ok := users[login]; ok {
		return nil, errors.New("User already exists " + login)
	}
	user := User{
		uuid:          uuid.New().ID(),
		login:         login,
		password_hash: password,
		email:         email,
		name:          name,
	}

	users[login] = user
	return &user, nil
}

func Auth(login string, password string) (*User, error) {
	user, ok := users[login]
	if !ok {
		return nil, errors.New("User not exists " + login)
	}
	if user.password_hash != password {
		return nil, errors.New("Wrong password " + login)
	}

	return &user, nil
}

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

func main() {
	users = make(map[string]User)
	sessions = make(map[string]Session)

	r := mux.NewRouter()

	r.HandleFunc("/api/auth", SessionMiddleware(HandleAuth)).Methods("POST")
	r.HandleFunc("/api/register", SessionMiddleware(HandleRegister)).Methods("POST")
	fs := http.FileServer(http.Dir("static/"))

	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
