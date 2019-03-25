package handlers

import (
	"encoding/json"
	_ "fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	webJson "../json_structs"
	"../models"
)

func HandleLogin(w http.ResponseWriter, r *http.Request, session *models.Session) {
	userData := &webJson.UsrRequest{}

	err := getRequest(userData, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := webJson.Response{
		Type: "log",
	}

	if session.User != nil {
		response.Status = "success"
	} else {
		user, err := models.Auth(userData.Login, userData.Password)
		if err != nil {
			wrong := err.Error()
			response.Status = "error"
			response.Payload = webJson.ErrorPayload{
				Message: "incorrect " + wrong,
				Field:   wrong,
			}
		} else {
			session.User = user
			response.Status = "success"
		}

		if response.Status == "success" {
			response.Payload = webJson.UserDataPayload{
				Login:      user.Login,
				Email:      user.Email,
				Name:       user.Name,
				AvatarPath: user.Avatar,
				Score:      user.Score,
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
func HandleRegister(w http.ResponseWriter, r *http.Request, session *models.Session) {
	userData := &webJson.UsrRequest{}

	err := getRequest(userData, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := webJson.Response{
		Type: "reg",
	}

	user, err := models.NewUser(userData.Login, userData.Password, userData.Email, userData.Name)
	if err == nil {
		session.User = user
		response.Status = "success"
		response.Payload = webJson.UserDataPayload{
			Login:      user.Login,
			Email:      user.Email,
			Name:       user.Name,
			AvatarPath: user.Avatar,
			Score:      user.Score,
		}
	} else {
		response.Status = "error"
		if err.Error() == "user already exists" {
			response.Payload = webJson.ErrorPayload{
				Message: err.Error(),
				Field:   "login",
			}
		} else {
			response.Payload = webJson.ErrorPayload{
				Message: "missing " + err.Error(),
				Field:   err.Error(),
			}
		}
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func HandleAvatarUpload(w http.ResponseWriter, r *http.Request, session *models.Session) {
	user := session.User
	response := webJson.Response{
		Type: "usinfo",
	}
	r.ParseMultipartForm(2 << 21) // 2 mb

	rFile, handler, err := r.FormFile("avatar")
	if err != nil {
		response.Status = "error"
		response.Payload = webJson.ErrorPayload{
			Message: "Wrong request",
			Field:   "avatar",
		}
	} else {
		defer rFile.Close()
		//fmt.Fprintf(w, "%v", handler.Header)
		filename := filepath.Join(filepath.Join("..", "2019_1_DeathPacito_front",
			"media", "avatar",
			uuid.New().String()+handler.Filename))

		wFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			response.Status = "error"
			response.Payload = webJson.ErrorPayload{
				Message: err.Error(),
				Field:   "avatar",
			}
		} else {
			defer wFile.Close()
			io.Copy(wFile, rFile)

			user.Avatar = filename
			err = user.Save()

			response.Status = "success"
			response.Payload = webJson.UserDataPayload{
				Login:      user.Login,
				Email:      user.Email,
				Name:       user.Name,
				AvatarPath: user.Avatar,
				Score:      user.Score,
			}
		}
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func HandleGetUsers(w http.ResponseWriter, r *http.Request, session *models.Session) {
	response := webJson.Response{
		Type: "uslist",
	}
	page, err := strconv.Atoi(mux.Vars(r)["page"])

	if err != nil {
		response.Status = "error"
		response.Payload = webJson.ErrorPayload{
			Message: "Wrong request",
		}
	} else {
		userSlice, err := models.GetUsers(10, page)

		if err != nil {
			response.Status = "error"
			response.Payload = webJson.ErrorPayload{
				Message: err.Error(),
			}
		} else {
			response.Status = "success"

			dataSlice := make([]webJson.UserDataPayload, 0, len(userSlice))
			for _, user := range userSlice {
				dataSlice = append(dataSlice, webJson.UserDataPayload{
					Name:  user.Name,
					Score: user.Score,
				})
			}
			count, _ := models.GetUserCount()
			response.Payload = webJson.UsersPayload{
				Users: dataSlice,
				Count: count,
			}
		}
	}
	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func HandleGetUserData(w http.ResponseWriter, r *http.Request, session *models.Session) {
	user := session.User
	response := webJson.Response{
		Type:   "usinfo",
		Status: "success",
	}

	response.Payload = webJson.UserDataPayload{
		Login:      user.Login,
		Email:      user.Email,
		Name:       user.Name,
		AvatarPath: user.Avatar,
		Score:      user.Score,
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func HandleUpdateUser(w http.ResponseWriter, r *http.Request, session *models.Session) {
	userData := &webJson.UsrRequest{}

	err := getRequest(userData, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := session.User

	if userData.Name != "" {
		user.Name = userData.Name
	}

	if userData.Password != "" {
		user.PasswordHash = userData.Password
	}

	user.Save()

	response := webJson.Response{
		Type:   "usinfo",
		Status: "success",
	}

	response.Payload = webJson.UserDataPayload{
		Login:      user.Login,
		Email:      user.Email,
		Name:       user.Name,
		AvatarPath: user.Avatar,
		Score:      user.Score,
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
