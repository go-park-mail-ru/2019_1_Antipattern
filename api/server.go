package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	pb "../api_struct"
	"../providers/auth"
	"./handlers"
	"./middleware"
	"./models"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type Config struct {
	APIPort			string							`json:"api_port"`
	AuthPort		string							`json:"auth_port"`
	APIPrefix		string							`json:"api_prefix"`
	AuthRoute		string							`json:"auth_route"`
	RegRoute		string							`json:"reg_route"`
	AvatarRoute		string							`json:"avatar_route"`
	ProfileRoute	string							`json:"profile_route"`
	LeaderRoute		string							`json:"leader_route"`
	UserRoute		string							`json:"user_route"`
	LoginRoute		string							`json:"login_route"`
	ServerAddress	string							`json:"server_address"`
	AuthDomain		string							`json:"auth_domain"`
}

var (
	config = &Config{}
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
		ServerAddress: config.ServerAddress + config.AuthPort,
		Secure:        false,
		AuthDomain:    config.AuthDomain,
	}
	r.HandleFunc(config.APIPrefix + config.AuthRoute, HandlerWrapperUnauthorized(handlers.HandleLogin, authProvider)).Methods("POST")
	r.HandleFunc(config.APIPrefix + config.RegRoute, HandlerWrapperUnauthorized(handlers.HandleRegister, authProvider)).Methods("POST")
	r.HandleFunc(config.APIPrefix + config.AvatarRoute, middleware.AuthMiddleware(handlers.HandleAvatarUpload, authProvider)).Methods("POST")
	r.HandleFunc(config.APIPrefix + config.ProfileRoute, middleware.AuthMiddleware(handlers.HandleUpdateUser, authProvider)).Methods("PUT")
	r.HandleFunc(config.APIPrefix + config.ProfileRoute, middleware.AuthMiddleware(handlers.HandleGetUserData, authProvider)).Methods("GET")
	r.HandleFunc(config.APIPrefix + config.LeaderRoute, handlers.HandleGetUsers).Methods("GET")
	r.HandleFunc(config.APIPrefix + config.UserRoute, handlers.HandleGetUserByID).Methods("GET")
	r.HandleFunc(config.APIPrefix + config.LoginRoute, middleware.AuthMiddleware(HandlerWrapperAuthorized(handlers.HandleLogout, authProvider), authProvider)).Methods("DELETE")
	return r
}
func main() {
	configBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Readn't: %v", err)
	}

	err = json.Unmarshal(configBytes, config)
	if err != nil {
		log.Fatalf("Unmarshalln't: %v", err)
	}

	models.InitModels(false)
	listener, err := net.Listen("tcp", config.AuthPort)
	log.Printf("API grpc server listening on " + config.AuthPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAPIServer(s, &server{})
	go s.Serve(listener)

	defer models.FinalizeModels()
	log.Fatal(http.ListenAndServe(config.APIPort, middleware.PanicMiddleware(NewRouter())))
}
