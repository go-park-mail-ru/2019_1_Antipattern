package main

import (
	"log"
	"net/http"
	"path"

	"./handlers"
	"./middleware"
	"./models"

	gHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	allowOrigins := gHandlers.AllowedOrigins([]string{"http://kpacubo.xyz", "http://api.kpacubo.xyz"})
	allowHeaders := gHandlers.AllowedHeaders([]string{"X-Requested-With"})
	allowMethods := gHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	r := mux.NewRouter()

	r.HandleFunc("/api/auth", middleware.SessionMiddleware(handlers.HandleLogin, false)).Methods("POST")
	r.HandleFunc("/api/register", middleware.SessionMiddleware(handlers.HandleRegister, false)).Methods("POST")
	r.HandleFunc("/api/upload_avatar", middleware.SessionMiddleware(handlers.HandleAvatarUpload, true)).Methods("POST")
	r.HandleFunc("/api/profile", middleware.SessionMiddleware(handlers.HandleUpdateUser, true)).Methods("PUT")
	r.HandleFunc("/api/profile", middleware.SessionMiddleware(handlers.HandleGetUserData, true)).Methods("GET")
	r.HandleFunc("/api/leaderboard/{page:[0-9]+}", middleware.SessionMiddleware(handlers.HandleGetUsers, false)).Methods("GET")

	staticServer := http.FileServer(http.Dir(
		path.Join("..", "2019_1_DeathPacito_front", "public")))
	mediaServer := http.FileServer(http.Dir("media/"))

	r.PathPrefix("/media").Handler(http.StripPrefix("/media/", mediaServer))
	r.PathPrefix("/public").Handler(http.StripPrefix("/public/", staticServer))

	r.HandleFunc("/", middleware.SessionMiddleware(func(w http.ResponseWriter, r *http.Request, s *models.Session) {
		http.ServeFile(w, r, path.Join(
			"..", "2019_1_DeathPacito_front",
			"public", "index.html"))
	}, false))

	gHandlers.AllowCredentials()
	return gHandlers.CORS(allowOrigins, allowHeaders, allowMethods)(r)
}
func main() {
	models.InitModels()

	log.Fatal(http.ListenAndServe(":8080", NewRouter()))
}
