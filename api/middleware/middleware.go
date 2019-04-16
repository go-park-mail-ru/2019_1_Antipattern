package middleware

import (
	"fmt"
	"net/http"

	"../models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

/*func SessionMiddleware(next func(http.ResponseWriter, *http.Request, *models.Session), authRequiered bool) http.HandlerFunc {
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

		if authRequiered && session.User == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if authRequiered {
			fmt.Println(session.User.Login)
		}
		next(w, r, &session)
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
	}
}
*/
