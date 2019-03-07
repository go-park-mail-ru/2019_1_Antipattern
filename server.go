package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	users = make(map[string]User)
	uuidUserIndex = make(map[uint32]string)
	sessions = make(map[string]Session)

	r := mux.NewRouter()

	r.HandleFunc("/api/auth", SessionMiddleware(HandleLogin, false)).Methods("POST")
	r.HandleFunc("/api/register", SessionMiddleware(HandleRegister, false)).Methods("POST")
	r.HandleFunc("/api/upload_avatar", SessionMiddleware(HandleAvatarUpload, true)).Methods("POST")
	r.HandleFunc("/api/profile", SessionMiddleware(HandleUpdateUser, true)).Methods("PUT")
	r.HandleFunc("/api/profile", SessionMiddleware(HandleGetUserData, true)).Methods("GET")
	r.HandleFunc("/api/leaderbord/{page:[0-9]+}", SessionMiddleware(HandleGetUsers, true)).Methods("GET")
	//staticServer := http.FileServer(http.Dir("../2019_1_DeathPacito_front/public"))
	staticServer := http.FileServer(http.Dir("../2019_1_DeathPacito_front/public/"))
	mediaServer := http.FileServer(http.Dir("media/"))
	//r.PathPrefix("/public").Handler(http.StripPrefix("/public/", staticServer))
	r.PathPrefix("/media").Handler(http.StripPrefix("/media/", mediaServer))
	r.PathPrefix("/public").Handler(http.StripPrefix("/public/", staticServer))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		http.ServeFile(w, r, "../2019_1_DeathPacito-front/public/index.html")
	})

	log.Fatal(http.ListenAndServe(":80", r))
}
