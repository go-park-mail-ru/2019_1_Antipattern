package main

import (
	"context"
	"log"
	"net"
	"net/http"

	pb "../api_struct"
	"../auth"
	"./handlers"
	"./middleware"
	"./models"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) GetUsers(ctx context.Context, request *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, err := models.GetUsersByIds(request.Uid)
	if err != nil {
		return nil, err
	}
	response := pb.GetUsersResponse{}
	for _, user := range users {
		data := pb.GetUsersResponse_UserData{Uid: user.Uuid.Hex(), Login: user.Login, Avatar: user.Avatar}
		response.Data = append(response.Data, &data)
	}
	return &response, nil
}

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
	listener, err := net.Listen("tcp", ":8081")
	log.Printf("API grpc server listening on 8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAPIServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	defer models.FinalizeModels()
	log.Fatal(http.ListenAndServe(":8080", middleware.PanicMiddleware(NewRouter())))
}
