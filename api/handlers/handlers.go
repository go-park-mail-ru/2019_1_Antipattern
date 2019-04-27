package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"

	webJson "../json_structs"
	"../models"
)

func setJWT(w http.ResponseWriter, user *models.User) error {
	secret := []byte("secret")
	uidHex := ""
	if user != nil {
		uidHex = user.Uuid.Hex()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uidHex,
		"sid": uuid.New().String(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("Failed to create token")
	}
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Domain:   ".kpacubo.xyz",
	}
	http.SetCookie(w, cookie)
	return nil
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	userData := &webJson.UsrRequest{}

	err := getRequest(userData, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := webJson.Response{
		Type: "log",
	}

	user, err := models.Auth(userData.Login, userData.Password)
	if err != nil {
		wrong := err.Error()
		response.Status = "error"
		response.Payload = webJson.ErrorPayload{
			Message: "incorrect " + wrong,
			Field:   wrong,
		}
	} else {
		setJWT(w, user)
		response.Status = "success"
	}

	if response.Status == "success" {
		response.Payload = webJson.UserDataPayload{
			Login:      user.Login,
			Email:      user.Email,
			AvatarPath: user.Avatar,
			Score:      user.Score,
		}
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

// HandleRegister handle registration api request
// request must contain post form:
// 	login, password, email
// Writes status json to response
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	userData := &webJson.UsrRequest{}

	err := getRequest(userData, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := webJson.Response{
		Type: "reg",
	}

	user, err := models.NewUser(userData.Login, userData.Password, userData.Email)
	if err == nil {
		setJWT(w, user)
		response.Status = "success"
		response.Payload = webJson.UserDataPayload{
			Login:      user.Login,
			Email:      user.Email,
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

func HandleAvatarUpload(w http.ResponseWriter, r *http.Request, user *models.User) {
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
		avatarName := uuid.New().String() + handler.Filename
		filename := filepath.Join(filepath.Join("/", "opt", "media", "avatar",
			avatarName))

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

			user.Avatar = "/media/avatar/" + avatarName
			err = user.Save()

			response.Status = "success"
			response.Payload = webJson.UserDataPayload{
				Login:      user.Login,
				Email:      user.Email,
				AvatarPath: user.Avatar,
				Score:      user.Score,
			}
		}
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
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
					Login: user.Login,
					Score: user.Score,
				})
			}
			count, _ := models.GetUserCount()
			response.Payload = webJson.UsersPayload{
				Users: dataSlice,
				Count: int(count),
			}
		}
	}
	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func HandleGetUserData(w http.ResponseWriter, r *http.Request, user *models.User) {
	response := webJson.Response{
		Type:   "usinfo",
		Status: "success",
	}

	response.Payload = webJson.UserDataPayload{
		Login:      user.Login,
		Email:      user.Email,
		AvatarPath: user.Avatar,
		Score:      user.Score,
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)
}

func HandleUpdateUser(w http.ResponseWriter, r *http.Request, user *models.User) {
	userData := &webJson.UsrRequest{}

	err := getRequest(userData, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userData.Login != "" {
		user.Login = userData.Login
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

func HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := webJson.Response{
		Type:   "usinfo",
		Status: "success",
	}
	objectID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := models.GetUser(objectID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response.Payload = webJson.UserDataPayload{
		Login:      user.Login,
		Email:      "",
		AvatarPath: user.Avatar,
		Score:      user.Score,
	}

	byteResponse, _ := response.MarshalJSON()
	w.Write(byteResponse)

}

func HandleLogout(w http.ResponseWriter, r *http.Request, user *models.User) {
	setJWT(w, nil)
}
