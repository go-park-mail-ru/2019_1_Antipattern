package auth

import "net/http"

// Provider is an interface for auth providers
type Provider interface {
	SetAuthCookie(w http.ResponseWriter, r *http.Request, uid string) error
	AuthMiddleware(next func(w http.ResponseWriter, r *http.Request, uid string)) http.HandlerFunc
	GetUUID(r *http.Request) (string, error)
	DeleteUserSession(w http.ResponseWriter, r *http.Request) error
	DeleteAllUserSessions(w http.ResponseWriter, r *http.Request) error
}
