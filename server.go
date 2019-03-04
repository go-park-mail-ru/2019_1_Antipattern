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

	r.HandleFunc("/api/auth", SessionMiddleware(HandleLogin)).Methods("POST")
	r.HandleFunc("/api/register", SessionMiddleware(HandleRegister)).Methods("POST")
	fs := http.FileServer(http.Dir("static/"))

	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", r))
}
