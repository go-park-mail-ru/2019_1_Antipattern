package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
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
			AvatarPath: user.avatar,
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

func HandleAvatarUpload(w http.ResponseWriter, r *http.Request, session *Session) {
	if session.user == nil {
		// TODO: write unauthorized to response
		return
	}
	r.ParseMultipartForm(2 << 21) // 2 mb
	rFile, handler, err := r.FormFile("avatar")
	if err != nil {
		// TODO: write error to response
		fmt.Println(err)
		return
	}
	defer rFile.Close()
	//fmt.Fprintf(w, "%v", handler.Header)
	filename := filepath.Join(filepath.Join(".", "media", "avatar",
		uuid.New().String()+handler.Filename))

	wFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		// TODO: write error to response
		fmt.Println(err)
		return
	}
	defer wFile.Close()
	io.Copy(wFile, rFile)

	session.user.avatar = filename
	err = session.user.Save()
	if err == nil {
		// TODO: write error to response
	}
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
