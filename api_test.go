package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// User register

func CheckSessionSetCookie(t *testing.T, user User, w *httptest.ResponseRecorder) {
	cookiesString := w.HeaderMap.Get("Set-Cookie")
	if cookiesString == "" {
		t.Errorf("Cookies not set")
		return
	}
	header := http.Header{}
	header.Add("Cookie", cookiesString)
	requestCooies := http.Request{Header: header}
	sessionID, err := requestCooies.Cookie("sid")
	if err != nil {
		t.Errorf("Session cookie not set")
		return
	}

	session, err := GetSession(sessionID.Value)
	if err != nil {
		t.Errorf("Can't get session!\n%s", err.Error())
		return
	}
	if session.user.uuid != user.uuid {
		t.Errorf("Session uuid is wrong.\nExpected:%s\nGot:%s", string(user.uuid), string(session.user.uuid))
		return
	}
}
func TestRegister(t *testing.T) {
	InitModels()

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
	expectedBody := `{"type":"reg","status":"success","payload":{"login":"user_login","email":"death.pa_cito@mail.yandex.ru","name":"Gamer #23 @790-_%"}}`

	w := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(w, r)

	result, _ := ioutil.ReadAll(w.Body)
	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result)
		return
	}
	newUser, _ := GetUserByLogin("user_login")

	if newUser.login != "user_login" || newUser.passwordHash != "qweqwe234234&62342=" ||
		newUser.email != "death.pa_cito@mail.yandex.ru" || newUser.name != "Gamer #23 @790-_%" {
		t.Errorf("Wrong user in db")
		return
	}
	CheckSessionSetCookie(t, *newUser, w)
}

func TestLogin(t *testing.T) {
	InitModels()
	user, err := NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru", "kek")
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
	expectedBody := `{"type":"log","status":"success","payload":{"login":"user_login","email":"death.pa_cito@mail.yandex.ru","name":"kek"}}`

	w := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(w, r)

	result, _ := ioutil.ReadAll(w.Body)
	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result)
		return
	}
	CheckSessionSetCookie(t, *user, w)
}
