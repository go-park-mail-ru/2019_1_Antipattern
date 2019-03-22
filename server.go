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

	r.HandleFunc("/api/auth", middleware.SessionMiddleware(handlers.HandleLogin, false)).Methods("POST")
	r.HandleFunc("/api/register", middleware.SessionMiddleware(handlers.HandleRegister, false)).Methods("POST")
	r.HandleFunc("/api/upload_avatar", middleware.SessionMiddleware(handlers.HandleAvatarUpload, true)).Methods("POST")
	r.HandleFunc("/api/profile", middleware.SessionMiddleware(handlers.HandleUpdateUser, true)).Methods("PUT")
	r.HandleFunc("/api/profile", middleware.SessionMiddleware(handlers.HandleGetUserData, true)).Methods("GET")
	r.HandleFunc("/api/leaderboard/{page:[0-9]+}", middleware.SessionMiddleware(handlers.HandleGetUsers, false)).Methods("GET")

	//staticServer := http.FileServer(http.Dir(
	//	path.Join("..", "2019_1_DeathPacito_front", "public")))
	//mediaServer := http.FileServer(http.Dir("media/"))

	//r.PathPrefix("/media").Handler(http.StripPrefix("/media/", mediaServer))
	//r.PathPrefix("/public").Handler(http.StripPrefix("/public/", staticServer))

	//r.HandleFunc("/", middleware.SessionMiddleware(func(w http.ResponseWriter, r *http.Request, s *models.Session) {
	//		http.ServeFile(w, r, path.Join(
	//			"..", "2019_1_DeathPacito_front",
	//			"public", "index.html"))
	//	}, false))

	return r
}
func main() {
	models.InitModels()

	log.Fatal(http.ListenAndServe(":8080", NewRouter()))
}
