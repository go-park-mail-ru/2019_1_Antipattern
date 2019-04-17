package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"./models"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// User register

func CheckSessionSetCookie(t *testing.T, user models.User, w *httptest.ResponseRecorder) {
	secret := []byte("secret")
	cookiesString := w.HeaderMap.Get("Set-Cookie")
	if cookiesString == "" {
		t.Errorf("Cookies not set")
		return
	}
	header := http.Header{}
	header.Add("Cookie", cookiesString)
	requestCooies := http.Request{Header: header}
	tokenString, err := requestCooies.Cookie("token")

	if err != nil {
		t.Errorf("Session cookie not set")
		return
	}
	token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("Can't get session!\n%s", err.Error())

		}
		return secret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// TODO: check type assertion
		uid := claims["uid"].(string)

		if uid != user.Uuid.Hex() {
			t.Errorf("Session uuid is wrong.\nExpected:%s\nGot:%s", user.Uuid.Hex(), uid)
			return
		}

	} else {
		t.Errorf("Can't get session!\n%s", err.Error())
		return
	}

	/*session, err := models.GetSession(sessionID.Value)
	if err != nil {
		t.Errorf("Can't get session!\n%s", err.Error())
		return
	}
	if session.User.Uuid != user.Uuid {
		t.Errorf("Session uuid is wrong.\nExpected:%s\nGot:%s", user.Uuid.String(), session.User.Uuid.String())
		return
	}*/
}

func SendApiQuery(request *http.Request, expectedBody string) (*httptest.ResponseRecorder, error) {
	response := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(response, request)
	result, _ := ioutil.ReadAll(response.Body)

	if strings.TrimSpace(string(result)) != expectedBody {
		return response, errors.New(fmt.Sprintf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result))
	}
	return response, nil
}
func TestRegister(t *testing.T) {
	models.InitModels(true)

	body := strings.NewReader(`{
		"login":"user_login",
		"password" : 
		"qweqwe234234&62342=",
		"email": 
		"death.pa_cito@mail.yandex.ru"}`)

	r, err := http.NewRequest("POST", "http://localhost/api/register", body)
	if err != nil {
		t.Fatal("Can't initialize")
		return
	}
	expectedBody := `{"type":"reg","status":"success","payload":{"login":"user_login","email":"death.pa_cito@mail.yandex.ru","score":20}}`

	response, err := SendApiQuery(r, expectedBody)
	if err != nil {
		t.Errorf(err.Error())
		//return
	}
	newUser, _ := models.GetUserByLogin("user_login")

	if newUser.Login != "user_login" || newUser.PasswordHash != "qweqwe234234&62342=" ||
		newUser.Email != "death.pa_cito@mail.yandex.ru" {
		t.Errorf("Wrong user in db %+v", newUser)
		return
	}
	CheckSessionSetCookie(t, *newUser, response)
}
func TestRegisterAlreadyRegistered(t *testing.T) {
	models.InitModels(true)
	expectedBody := `{"type":"reg","status":"error","payload":{"message":"user already exists","field":"login"}}`

	_, err := models.NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru")
	if err != nil {
		t.Fatal("Can't create user")
		return
	}
	body := strings.NewReader(`{
		"login":"user_login",
		"password" : 
		"qweqwe234234&62342=",
		"email": 
		"death.pa_cito@mail.yandex.ru",
		"name": "Gamer #23 @790-_%" }`)

	r, err := http.NewRequest("POST", "http://localhost/api/register", body)
	if err != nil {
		t.Fatal("Can't initialize")
		return
	}

	_, err = SendApiQuery(r, expectedBody)
	if err != nil {
		t.Errorf(err.Error())
	}
}
func TestLogin(t *testing.T) {
	expectedBody := `{"type":"log","status":"success","payload":{"login":"user_login","email":"death.pa_cito@mail.yandex.ru","score":20}}`
	models.InitModels(true)
	user, err := models.NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru")
	if err != nil {
		t.Fatal("Can't create user")
		return
	}
	body := strings.NewReader(`{
		"login":"user_login",
		"password" : "1235689"}`)

	r, err := http.NewRequest("POST", "http://localhost/api/auth", body)
	if err != nil {
		t.Fatal("Can't initialize")
		return
	}

	response, err := SendApiQuery(r, expectedBody)
	if err != nil {
		t.Errorf(err.Error())
	}

	CheckSessionSetCookie(t, *user, response)
}

func TestLoginWrongPassword(t *testing.T) {
	models.InitModels(true)
	expectedBody := `{"type":"log","status":"error","payload":{"message":"incorrect password","field":"password"}}`
	_, err := models.NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru")
	if err != nil {
		t.Fatal("Can't create user")
		return
	}
	body := strings.NewReader(`{
		"login":"user_login",
		"password" : "asdasd"}`)

	r, err := http.NewRequest("POST", "http://localhost/api/auth", body)
	if err != nil {
		t.Fatal("Can't initialize")
		return
	}

	_, err = SendApiQuery(r, expectedBody)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func DisabledTestLoginWrongLogin(t *testing.T) {
	models.InitModels(true)
	expectedBody := `{"type":"log","status":"error","payload":{"message":"incorrect login","field":"login"}}`
	_, err := models.NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru")
	if err != nil {
		t.Fatal("Can't create user")
		return
	}
	body := strings.NewReader(`{
		"login":"lol",
		"password" : "1235689"}`)

	r, err := http.NewRequest("POST", "http://localhost/api/auth", body)
	if err != nil {
		t.Fatal("Can't initialize")
		return
	}

	_, err = SendApiQuery(r, expectedBody)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func FakeLoginAndAuth(request *http.Request) (*models.User, error) {
	user, err := models.NewUser("fake_user_login", "12345", "mail@mail.ru")
	if err != nil {
		return nil, err
	}
	hmacSampleSecret := []byte("secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.Uuid.Hex(),
		"sid": uuid.New().String(),
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return nil, errors.New("Cannot fake login!")
	}
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
	}
	request.AddCookie(cookie)
	return user, err
}

func TestGetProfile(t *testing.T) {
	models.InitModels(true)
	request, err := http.NewRequest("GET", "http://localhost/api/profile", nil)
	expectedBody := `{"type":"usinfo","status":"success","payload":{"login":"fake_user_login","email":"mail@mail.ru","score":20}}`
	_, err = FakeLoginAndAuth(request)

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	_, err = SendApiQuery(request, expectedBody)

	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestUpdateProfile(t *testing.T) {
	models.InitModels(true)
	body := strings.NewReader(`{
		"password" : "qweqwe234234&62342=",
		"name": "new name" }`)
	request, err := http.NewRequest("PUT", "http://localhost/api/profile", body)
	expectedBody := `{"type":"usinfo","status":"success","payload":{"login":"fake_user_login","email":"mail@mail.ru","score":20}}`

	user, err := FakeLoginAndAuth(request)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = SendApiQuery(request, expectedBody)
	if err != nil {
		t.Errorf(err.Error())
	}

	user, _ = models.GetUserByLogin("fake_user_login")

	if user.PasswordHash != "qweqwe234234&62342=" {
		t.Errorf("Wrong password hash\n Expected:qweqwe234234&62342=\nGot:%s", user.PasswordHash)
	}
}

func TestGetLeaderboard(t *testing.T) {
	models.InitModels(true)

	for i := 1; i <= 27; i++ {
		a := strconv.Itoa(i)
		models.NewUser("npc_"+a, "12345", "mail"+a+"@mail.ru")
	}
	request, err := http.NewRequest("GET", "http://localhost/api/leaderboard/1", nil)
	expectedBody := `{"type":"uslist","status":"success","payload":{"users":[{"login":"fake_user_login","score":20},{"login":"npc_1","score":20},{"login":"npc_10","score":20},{"login":"npc_11","score":20},{"login":"npc_12","score":20},{"login":"npc_13","score":20},{"login":"npc_14","score":20},{"login":"npc_15","score":20},{"login":"npc_16","score":20},{"login":"npc_17","score":20}],"count":28}}`
	_, err = FakeLoginAndAuth(request)

	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = SendApiQuery(request, expectedBody)

	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestGetLeaderboardTooBigPage(t *testing.T) {
	models.InitModels(true)
	expectedBody := `{"type":"uslist","status":"error","payload":{"message":"not enough users"}}`

	for i := 1; i <= 27; i++ {
		models.NewUser("npc_"+string(i), "12345", "mail"+string(i)+"@mail.ru")
	}
	request, err := http.NewRequest("GET", "http://localhost/api/leaderboard/31", nil)

	_, err = FakeLoginAndAuth(request)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = SendApiQuery(request, expectedBody)
	if err != nil {
		t.Errorf(err.Error())
	}

}
