package main

import (
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/auth", SessionMiddleware(HandleLogin, false)).Methods("POST")						// check, но изменить ошибки
	r.HandleFunc("/api/register", SessionMiddleware(HandleRegister, false)).Methods("POST")					// принимает неполные запросыFFF
	r.HandleFunc("/api/upload_avatar", SessionMiddleware(HandleAvatarUpload, true)).Methods("POST")			//
	r.HandleFunc("/api/profile", SessionMiddleware(HandleUpdateUser, true)).Methods("PUT")					//
	r.HandleFunc("/api/profile", SessionMiddleware(HandleGetUserData, true)).Methods("GET")					// хз вроде норм
	r.HandleFunc("/api/leaderboard/{page:[0-9]+}", SessionMiddleware(HandleGetUsers, true)).Methods("GET")	// -

	staticServer := http.FileServer(http.Dir(
		path.Join("..", "2019_1_DeathPacito_front", "public")))
	mediaServer := http.FileServer(http.Dir("media/"))

	r.PathPrefix("/media").Handler(http.StripPrefix("/media/", mediaServer))
	r.PathPrefix("/public").Handler(http.StripPrefix("/public/", staticServer))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(
			"..", "2019_1_DeathPacito_front",
			"public", "index.html"))
	})
	return r
}
func main() {
	InitModels()

	log.Fatal(http.ListenAndServe(":80", NewRouter()))
}
