package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func HandleLogin(w http.ResponseWriter, r *http.Request, session *Session) {
	userData := &UsrRequest{}
	err := getRequest(userData, r)
	if err != nil {
		fmt.Printf("An error occured: %v\nRequest: %v", err, userData)
		return
	}

	response := Response{
		Type:    "log",
		Payload: nil,
	}

	if session.user != nil {
		response.Status = "success"
	} else {
		user, err := Auth(userData.Login, userData.Password)
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
				AvatarPath: user.avatar,
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
	userData := &UsrRequest{}
	err := getRequest(userData, r)
	if err != nil {
		fmt.Printf("An error occured: %v\nRequest: %v", err, userData)
		return
	}

	response := Response{
		Type: "reg",
	}
	
	user, err := NewUser(userData.Login, userData.Password, userData.Email, userData.Name)
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
	r.ParseMultipartForm(2 << 21) // 2 mb
	rFile, handler, err := r.FormFile("avatar")
	if err != nil {
		// TODO: write error to response
		fmt.Println(err)
		return
	}
	defer rFile.Close()
	//fmt.Fprintf(w, "%v", handler.Header)
	filename := filepath.Join(filepath.Join("media", "avatar",
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

func HandleGetUsers(w http.ResponseWriter, r *http.Request, session *Session) {
	request := &LeaderboardRequest{}
	err := getRequest(request, r)
	if err != nil {
		fmt.Printf("An error occured: %v\nRequest: %v", err, request)
		return
	}

	response := Response{
		Type: "uslist",
	}

	userSlice, err := GetUsers(request.Count, request.Page)
	if err != nil {
		response.Status = "error"
		response.Payload = ErrorPayload{
			Message: "Not enough users",
		}
	} else {
		response.Status = "success"
		response.Payload = UsersPayload{
			Users: userSlice,
		}
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func HandleGetUserData(w http.ResponseWriter, r *http.Request, session *Session) {
	user := session.user
	response := Response{
		Type:   "usinfo",
		Status: "success",
	}

	response.Payload = UserDataPayload{
		Login:      user.login,
		Email:      user.email,
		Name:       user.name,
		AvatarPath: user.name,
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func HandleUpdateUser(w http.ResponseWriter, r *http.Request, session *Session) {
	userData := &UsrRequest{}
	err := getRequest(userData, r)
	if err != nil {
		fmt.Printf("An error occured: %v\nRequest: %v", err, userData)
		return
	}

	user := session.user

	if userData.Name != "" {
		user.name = userData.Name
	}

	if userData.Password != "" {
		user.name = userData.Password
	}

	user.Save()

	response := Response{
		Type:   "usinfo",
		Status: "success",
	}

	response.Payload = UserDataPayload{
		Login:      user.login,
		Email:      user.email,
		Name:       user.name,
		AvatarPath: user.avatar,
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func getRequest(marshaler json.Unmarshaler, r *http.Request) error {
	body := r.Body
	defer body.Close()
	byteBody, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	marshaler.UnmarshalJSON(byteBody)

	if err != nil {
		return err
	}
	return nil
}
