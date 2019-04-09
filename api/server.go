package main

import (
	"log"
	"net/http"

	"./handlers"
	"./middleware"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {

	r := mux.NewRouter()

	r.HandleFunc("/api/auth", middleware.SessionMiddleware(handlers.HandleLogin, false)).Methods("POST")
	r.HandleFunc("/api/register", middleware.SessionMiddleware(handlers.HandleRegister, false)).Methods("POST")
	r.HandleFunc("/api/upload_avatar", middleware.SessionMiddleware(handlers.HandleAvatarUpload, true)).Methods("POST")
	r.HandleFunc("/api/profile", middleware.SessionMiddleware(handlers.HandleUpdateUser, true)).Methods("PUT")
	r.HandleFunc("/api/profile", middleware.SessionMiddleware(handlers.HandleGetUserData, true)).Methods("GET")
	r.HandleFunc("/api/leaderboard/{page:[0-9]+}", middleware.SessionMiddleware(handlers.HandleGetUsers, false)).Methods("GET")
	r.HandleFunc("/api/login", middleware.SessionMiddleware(handlers.HandleLogout, true)).Methods("DELETE")
	return r
}
func main() {
	log.Fatal(http.ListenAndServe(":8080", NewRouter()))
}
