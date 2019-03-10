package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandleLogin(w http.ResponseWriter, r *http.Request, session *Session) {
	userData := &UsrRequest{}

	err := getRequest(userData, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := Response{
		Type: "log",
	}

	if session.user != nil {
		response.Status = "success"
	} else {
		user, err := Auth(userData.Login, userData.Password)
		if err != nil {
			wrong := err.Error()
			response.Status = "error"
			response.Payload = ErrorPayload{
				Message: "incorrect " + wrong,
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
				Score:      user.score,
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
		w.WriteHeader(http.StatusBadRequest)
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
			Score:      user.score,
		}
	} else {
		response.Status = "error"
		if err.Error() == "user already exists" {
			response.Payload = ErrorPayload{
				Message: err.Error(),
				Field:   "login",
			}
		} else {
			response.Payload = ErrorPayload{
				Message: "missing " + err.Error(),
				Field:   err.Error(),
			}
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
	response := Response{
		Type: "uslist",
	}
	page, err := strconv.Atoi(mux.Vars(r)["page"])

	if err != nil {
		response.Status = "error"
		response.Payload = ErrorPayload{
			Message: "Wrong request",
		}
	} else {
		userSlice, err := GetUsers(10, page)

		if err != nil {
			response.Status = "error"
			response.Payload = ErrorPayload{
				Message: err.Error(),
			}
		} else {
			response.Status = "success"

			dataSlice := make([]UserDataPayload, 0, len(userSlice))
			for _, user := range userSlice {
				dataSlice = append(dataSlice, UserDataPayload{
					Name:  user.name,
					Score: user.score,
				})
			}
			count, _ := GetUserCount()
			response.Payload = UsersPayload{
				Users: dataSlice,
				Count: count,
			}
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
		AvatarPath: user.avatar,
		Score:      user.score,
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func HandleUpdateUser(w http.ResponseWriter, r *http.Request, session *Session) {
	userData := &UsrRequest{}

	err := getRequest(userData, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := session.user

	if userData.Name != "" {
		user.name = userData.Name
	}

	if userData.Password != "" {
		user.passwordHash = userData.Password
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
		Score:      user.score,
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

	err = marshaler.UnmarshalJSON(byteBody)

	if err != nil {
		return err
	}
	return nil
}
