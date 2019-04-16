package middleware

import (
	"fmt"
	"net/http"

	"../models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PanicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				fmt.Println("Panic recovered")
			}
		}()
		h.ServeHTTP(w, r)
	})
}
func JWTMiddleware(next func(http.ResponseWriter, *http.Request, *models.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := []byte("secret")
		cookie, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			uid, ok := claims["uid"].(string)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if uid == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			objectID, err := primitive.ObjectIDFromHex(uid)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			user, err := models.GetUser(objectID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next(w, r, user)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

	}
}
