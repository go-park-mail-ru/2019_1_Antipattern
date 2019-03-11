package main

import (
	"log"
	"net/http"
	"path"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	allowOrigins := handlers.AllowedOrigins([]string{"http://kpacubo.xyz", "http://api.kpacubo.xyz", "kpacubo.xyz", "api.kpacubo.xyz"})
	allowHeaders := handlers.AllowedHeaders([]string{"X-Requested-With"})
	allowMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	r := mux.NewRouter()

	r.HandleFunc("/api/auth", SessionMiddleware(HandleLogin, false)).Methods("POST")                        // check, но изменить ошибки
	r.HandleFunc("/api/register", SessionMiddleware(HandleRegister, false)).Methods("POST")                 // принимает неполные запросыFFF
	r.HandleFunc("/api/upload_avatar", SessionMiddleware(HandleAvatarUpload, true)).Methods("POST")         //
	r.HandleFunc("/api/profile", SessionMiddleware(HandleUpdateUser, true)).Methods("PUT")                  //
	r.HandleFunc("/api/profile", SessionMiddleware(HandleGetUserData, true)).Methods("GET")                 // хз вроде норм
	r.HandleFunc("/api/leaderboard/{page:[0-9]+}", SessionMiddleware(HandleGetUsers, false)).Methods("GET") // -

	staticServer := http.FileServer(http.Dir(
		path.Join("..", "2019_1_DeathPacito_front", "public")))
	mediaServer := http.FileServer(http.Dir("media/"))

	r.PathPrefix("/media").Handler(http.StripPrefix("/media/", mediaServer))
	r.PathPrefix("/public").Handler(http.StripPrefix("/public/", staticServer))

	r.HandleFunc("/", SessionMiddleware(func(w http.ResponseWriter, r *http.Request, s *Session) {
		http.ServeFile(w, r, path.Join(
			"..", "2019_1_DeathPacito_front",
			"public", "index.html"))
	}, false))
	return handlers.CORS(allowOrigins, allowHeaders, allowMethods)(r)
}
func main() {
	InitModels()

	log.Fatal(http.ListenAndServe(":8080", NewRouter()))
}
