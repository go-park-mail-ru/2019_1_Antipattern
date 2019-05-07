package middleware

import (
	"fmt"
	"net/http"

	"../../auth"
	"../models"
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

func AuthMiddleware(next func(http.ResponseWriter, *http.Request, *models.User), authProvider auth.Provider) http.HandlerFunc {
	return authProvider.AuthMiddleware(func(w http.ResponseWriter, r *http.Request, uid string) {
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
	})
}
