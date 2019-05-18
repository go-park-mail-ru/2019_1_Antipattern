package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	pb "../../identity_struct"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

const (
	authCookieName    = "_auth_cookie"
	refreshCookieName = "_refresh_cookie"
)

// JWTProvider provides JWT based authoriztion
type JWTProvider struct {
	ServerAddress string
	Secure        bool
	AuthDomain    string
}

func (provider JWTProvider) SetAuthCookie(w http.ResponseWriter, r *http.Request, uid string) error {
	conn, err := grpc.Dial(provider.ServerAddress, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Printf("Can't connect to authority server: %v", err)
		return err
	}
	defer conn.Close()
	c := pb.NewIdentifierClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := c.IssueToken(ctx, &pb.IssueTokenRequest{Uid: uid})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	authCookie := &http.Cookie{
		Name:     authCookieName,
		Value:    response.AccessToken,
		HttpOnly: true,
		Domain:   provider.AuthDomain,
		Path:     "/",
	}
	refreshCookie := &http.Cookie{
		Name:     refreshCookieName,
		Value:    response.RefreshToken,
		HttpOnly: true,
		Domain:   provider.AuthDomain,
		Path:     "/",
	}
	http.SetCookie(w, authCookie)
	http.SetCookie(w, refreshCookie)
	return nil
}

// AuthMiddleware checks JWT cookies
func (provider JWTProvider) AuthMiddleware(next func(w http.ResponseWriter, r *http.Request, uid string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acccessCookie, err := r.Cookie(authCookieName)
		refreshCookie, err := r.Cookie(refreshCookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		conn, err := grpc.Dial(provider.ServerAddress, grpc.WithInsecure())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Can't connect to authority server: %v", err)
			return
		}
		defer conn.Close()
		c := pb.NewIdentifierClient(conn)

		response, err := c.ParseToken(context.Background(), &pb.ParseTokenRequest{
			AccessToken:  acccessCookie.Value,
			RefreshToken: refreshCookie.Value,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Can't parse auth token: %v", err)
			return
		}
		switch response.Status {
		case pb.ParseTokenStatusCode_SUCCESS:
			next(w, r, response.Uid)
		case pb.ParseTokenStatusCode_INVALID:
			w.WriteHeader(http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func (provider JWTProvider) GetUUID(r *http.Request) (string, error) {
	acccessCookie, err := r.Cookie(authCookieName)
	refreshCookie, err := r.Cookie(refreshCookieName)
	if err != nil {
		return "", err
	}
	conn, err := grpc.Dial(provider.ServerAddress, grpc.WithInsecure())

	if err != nil {
		log.Printf("Can't connect to authority server: %v", err)
		return "", err
	}
	defer conn.Close()
	c := pb.NewIdentifierClient(conn)

	response, err := c.ParseToken(context.Background(), &pb.ParseTokenRequest{
		AccessToken:  acccessCookie.Value,
		RefreshToken: refreshCookie.Value,
	})

	if err != nil {
		log.Printf("Can't parse auth token: %v", err)
		return "", err
	}
	switch response.Status {
	case pb.ParseTokenStatusCode_SUCCESS:
		return response.Uid, nil
	case pb.ParseTokenStatusCode_INVALID:
		return "", errors.New("unauthorized")
	default:
		return "", errors.New("unauthorized")
	}
}
func (provider JWTProvider) DeleteUserSession(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (provider JWTProvider) DeleteAllUserSessions(http.ResponseWriter, *http.Request) error {
	return nil
}
