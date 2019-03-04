package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type User struct {
	uuid          string
	login         string
	password_hash string
	email         string
	name          string
}

type Session struct {
	user User
}

var users map[string]User
var sessions map[string]Session

func createSession(user *User) string {
	session := Session{
		user: *user,
	}
	id := uuid.New().String()
	sessions[id] = session
	return id
}

func Register(login string, password string, email string, name string) (string, error) {

	if _, ok := users[login]; ok {
		return "", errors.New("User already exists " + login)
	}
	user := User{
		login:         login,
		password_hash: password,
		email:         email,
		name:          name,
	}

	users[login] = user
	id := createSession(&user)
	return id, nil
}

func Auth(login string, password string) (string, error) {
	user, ok := users[login]
	if !ok {
		return "", errors.New("User not exists " + login)
	}
	if user.password_hash != password {
		return "", errors.New("Wrong password " + login)
	}
	id := createSession(&user)
	return id, nil
}

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	cookie, err := r.Cookie("uuid")
	// TODO: In case if cookie is broken user must auth
	if err != nil {
		id, err := Auth(login, password)
		if err != nil {
			fmt.Fprintln(w, err.Error())
		} else {
			cookie = &http.Cookie{Name: "uuid", Value: id, HttpOnly: true}
			http.SetCookie(w, cookie)
			fmt.Fprintln(w, "Succesfull", id)
		}
	} else {
		fmt.Fprintln(w, "Already logged in")
	}

}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	email := r.FormValue("email")
	name := r.FormValue("name")

	id, err := Register(login, password, email, name)
	if err == nil {
		cookie := &http.Cookie{Name: "uuid", Value: id, HttpOnly: true}
		http.SetCookie(w, cookie)

		fmt.Fprintln(w, "User registered:", id)
	} else {
		fmt.Fprintln(w, err.Error())
	}

}

func main() {
	users = make(map[string]User)
	sessions = make(map[string]Session)

	r := mux.NewRouter()

	r.HandleFunc("/api/auth", HandleAuth).Methods("POST")
	r.HandleFunc("/api/register", HandleRegister).Methods("POST")
	fs := http.FileServer(http.Dir("static/"))

	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
