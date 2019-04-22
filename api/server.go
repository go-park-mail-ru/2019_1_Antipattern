package main

import (
	"log"
	"net/http"

	"./handlers"
	"./middleware"
	"./models"
	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {

	r := mux.NewRouter()

	r.HandleFunc("/api/auth", handlers.HandleLogin).Methods("POST")
	r.HandleFunc("/api/register", handlers.HandleRegister).Methods("POST")
	r.HandleFunc("/api/upload_avatar", middleware.JWTMiddleware(handlers.HandleAvatarUpload)).Methods("POST")
	r.HandleFunc("/api/profile", middleware.JWTMiddleware(handlers.HandleUpdateUser)).Methods("PUT")
	r.HandleFunc("/api/profile", middleware.JWTMiddleware(handlers.HandleGetUserData)).Methods("GET")
	r.HandleFunc("/api/leaderboard/{page:[0-9]+}", handlers.HandleGetUsers).Methods("GET")
	r.HandleFunc("/api/login", middleware.JWTMiddleware(handlers.HandleLogout)).Methods("DELETE")
	return r
}
func main() {
	models.InitModels(false)
	defer models.FinalizeModels()
	log.Fatal(http.ListenAndServe(":8080", middleware.PanicMiddleware(NewRouter())))
}
