package main

import (
	"fmt"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request, session *Session) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	response := Response{
		Type: "log",
		Payload: nil,
	}

	if session.user != nil {
		response.Status = "success"
	} else {
		user, err := Auth(login, password)
		if err != nil {
			wrong := err.Error()
			response.Status = "error"
			response.Payload = ErrorPayload{
				Message: "Incorrect" + wrong,
				Field: wrong,
			}
		} else {
			session.user = user
			response.Status = "success"
		}

		if response.Status == "success" {
			response.Payload = UserDataPayload{
				Login: user.login,
				Email: user.email,
				Name: user.name,
				AvatarPath: "fish.jpg",
			}
		}
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

// HandleRegister handle registration api request
// request must contain post form:
// 	login, password, email, name
// Writes status json to response
func HandleRegister(w http.ResponseWriter, r *http.Request, session *Session) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	email := r.FormValue("email")
	name := r.FormValue("name")
	response := Response{
		Type: "reg",
	}

	user, err := NewUser(login, password, email, name)
	if err == nil {
		session.user = user
		response.Status = "success"
		response.Payload = UserDataPayload{
			Login: user.login,
			Email: user.email,
			Name: user.name,
			AvatarPath: "fish.jpg",
		}
	} else {
		fmt.Fprintln(w, err.Error())
		response.Status = "error"
		response.Payload = ErrorPayload{
			Message: "User already exists",
			Field: "login",
		}
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}
