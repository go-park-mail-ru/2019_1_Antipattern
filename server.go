package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	users = make(map[string]User)
	sessions = make(map[string]Session)

	r := mux.NewRouter()

	r.HandleFunc("/api/auth", SessionMiddleware(HandleLogin, false)).Methods("POST")
	r.HandleFunc("/api/register", SessionMiddleware(HandleRegister, false)).Methods("POST")
	r.HandleFunc("/api/upload_avatar", SessionMiddleware(HandleAvatarUpload, true)).Methods("POST")
	r.HandleFunc("/api/profile", SessionMiddleware(HandleUpdateUser, true)).Methods("PUT")
	r.HandleFunc("/api/profile", SessionMiddleware(HandleGetUserData, true)).Methods("GET")
	r.HandleFunc("/api/leaderbord/{page:[0-9]+}", SessionMiddleware(HandleGetUsers, true)).Methods("GET")
	fs := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	log.Fatal(http.ListenAndServe(":8081", r))
}
