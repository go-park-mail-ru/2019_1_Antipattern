package main

import (
	"log"
	"net/http"

	"../auth"
	"./handlers"
	"./middleware"
	"./models"
	"github.com/gorilla/mux"
)

func HandlerWrapperUnauthorized(handler func(w http.ResponseWriter, r *http.Request, authProvider auth.Provider), authProvider auth.Provider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, authProvider)
	}
}
func HandlerWrapperAuthorized(handler func(w http.ResponseWriter, r *http.Request, user *models.User, authProvider auth.Provider), authProvider auth.Provider) func(w http.ResponseWriter, r *http.Request, user *models.User) {
	return func(w http.ResponseWriter, r *http.Request, user *models.User) {
		handler(w, r, user, authProvider)
	}
}
func NewRouter() http.Handler {

	r := mux.NewRouter()
	authProvider := auth.JWTProvider{
		ServerAddress: "identity_service:8081",
		Secure:        false,
		AuthDomain:    ".kpacubo.xyz",
	}
	r.HandleFunc("/api/auth", HandlerWrapperUnauthorized(handlers.HandleLogin, authProvider)).Methods("POST")
	r.HandleFunc("/api/register", HandlerWrapperUnauthorized(handlers.HandleRegister, authProvider)).Methods("POST")
	r.HandleFunc("/api/upload_avatar", middleware.AuthMiddleware(handlers.HandleAvatarUpload, authProvider)).Methods("POST")
	r.HandleFunc("/api/profile", middleware.AuthMiddleware(handlers.HandleUpdateUser, authProvider)).Methods("PUT")
	r.HandleFunc("/api/profile", middleware.AuthMiddleware(handlers.HandleGetUserData, authProvider)).Methods("GET")
	r.HandleFunc("/api/leaderboard/{page:[0-9]+}", handlers.HandleGetUsers).Methods("GET")
	r.HandleFunc("/api/user/{id:[0-9A-Fa-f]+}", handlers.HandleGetUserByID).Methods("GET")
	r.HandleFunc("/api/login", middleware.AuthMiddleware(HandlerWrapperAuthorized(handlers.HandleLogout, authProvider), authProvider)).Methods("DELETE")
	return r
}
func main() {
	models.InitModels(false)
	defer models.FinalizeModels()
	log.Fatal(http.ListenAndServe(":8080", middleware.PanicMiddleware(NewRouter())))
}
