package main

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"
	"time"

	pb "../identity_struct"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var secret = []byte("secret")

type server struct{}

func (s *server) ParseToken(ctx context.Context, request *pb.ParseTokenRequest) (*pb.ParseTokenResponse, error) {

	token, err := jwt.Parse(request.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New("Unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, errors.New("Invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid, ok := claims["uid"].(string)
		if !ok {
			return nil, errors.New("Invalid uid")
		}
		response := pb.ParseTokenResponse{Uid: uid}
		return &response, nil
	}
	return nil, errors.New("Invalid token")

}

func (s *server) IssueToken(ctx context.Context, request *pb.IssueTokenRequest) (*pb.IssueTokenResponse, error) {
	accessExpire := time.Now().AddDate(0, 0, 1).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      request.Uid,
		"token_id": uuid.New().String(),
		"expire":   strconv.FormatInt(accessExpire, 10),
	})

	accessStr, err := accessToken.SignedString(secret)
	if err != nil {
		return nil, err
	}
	response := pb.IssueTokenResponse{
		AccessToken:  accessStr,
		RefreshToken: "none",
	}
	return &response, nil
}

func (s *server) RevokeUserTokens(ctx context.Context, request *pb.RevokeUserTokensRequest) (*pb.RevokeTokenResponse, error) {
	resp := pb.RevokeTokenResponse{}
	return &resp, nil
}

func (s *server) RevokeToken(ctx context.Context, request *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	resp := pb.RevokeTokenResponse{}
	return &resp, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8081")
	log.Printf("Identity server listening on 8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterIdentifierServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
