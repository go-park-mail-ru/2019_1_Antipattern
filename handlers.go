package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request, session *Session) {
	request, err := getRequest(r)
	if err != nil {
		fmt.Printf("An error occured: %v\nRequest: %v", err, request)
		return
	}

	response := Response{
		Type:    "log",
		Payload: nil,
	}

	if session.user != nil {
		response.Status = "success"
	} else {
		user, err := Auth(request.Login, request.Password)
		if err != nil {
			wrong := err.Error()
			response.Status = "error"
			response.Payload = ErrorPayload{
				Message: "Incorrect" + wrong,
				Field:   wrong,
			}
		} else {
			session.user = user
			response.Status = "success"
		}

		if response.Status == "success" {
			response.Payload = UserDataPayload{
				Login:      user.login,
				Email:      user.email,
				Name:       user.name,
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
	request, err := getRequest(r)
	if err != nil {
		fmt.Printf("An error occured: %v\nRequest: %v", err, request)
		return
	}

	response := Response{
		Type: "reg",
	}

	user, err := NewUser(request.Login, request.Password, request.Email, request.Name)
	if err == nil {
		session.user = user
		response.Status = "success"
		response.Payload = UserDataPayload{
			Login:      user.login,
			Email:      user.email,
			Name:       user.name,
			AvatarPath: "fish.jpg",
		}
	} else {
		response.Status = "error"
		response.Payload = ErrorPayload{
			Message: "User already exists",
			Field:   "login",
		}
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func getRequest(r *http.Request) (*Request, error) {
	body := r.Body
	defer body.Close()

	byteBody, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	request := &Request{}
	err = request.UnmarshalJSON(byteBody)
	if err != nil {
		return nil, err
	}

	return request, nil
}