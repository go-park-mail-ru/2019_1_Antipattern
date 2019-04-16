package middleware

import (
	"fmt"
	"net/http"

	"../models"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SessionMiddleware(next func(http.ResponseWriter, *http.Request, *models.Session), authRequiered bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := []byte("secret")
		cookie, err := r.Cookie("token")
		if err != nil {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"uid": "",
				"sid": uuid.New().String(),
			})
			tokenString, err := token.SignedString(secret)
			if err != nil {
				// TODO: write error to response
			}
			cookie = &http.Cookie{
				Name:     "token",
				Value:    tokenString,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}
		session := models.Session{User: nil}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// TODO: check type assertion
			uid := claims["uid"].(string)
			session.Sid = claims["sid"].(string)
			if uid == "" {
				session.User = nil
			} else {
				objectID, err := primitive.ObjectIDFromHex(uid)
				user, err := models.GetUser(objectID)
				if err != nil {
					session.User = nil
					fmt.Println("Auth error!")
				} else {
					session.User = user
				}
			}
		} else {
			fmt.Println(err)
		}
		userHex := ""
		if session.User != nil {
			userHex = session.User.Uuid.Hex()
		}

		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"uid": userHex,
			"sid": uuid.New().String(),
		})

		tokenString, err := token.SignedString(secret)
		if err != nil {
			// TODO: write error to response
		}

		cookie.Value = tokenString
		if authRequiered && session.User == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next(w, r, &session)
	}
}
