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
func TestRegisterAlreadyRegistered(t *testing.T) {
	InitModels()
	_, err := NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru", "kek")
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
	expectedBody := `{"type":"reg","status":"error","payload":{"message":"User already exists","field":"login"}}`

	w := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(w, r)

	result, _ := ioutil.ReadAll(w.Body)
	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result)
		return
	}
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

func TestLoginWrongPassword(t *testing.T) {
	InitModels()
	_, err := NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru", "kek")
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
	expectedBody := `{"type":"log","status":"error","payload":{"message":"Incorrectpassword","field":"password"}}`

	w := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(w, r)

	result, _ := ioutil.ReadAll(w.Body)
	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result)
		return
	}
}

func TestLoginWrongLogin(t *testing.T) {
	InitModels()
	_, err := NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru", "kek")
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
	expectedBody := `{"type":"log","status":"error","payload":{"message":"Incorrectlogin","field":"login"}}`

	w := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(w, r)

	result, _ := ioutil.ReadAll(w.Body)
	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result)
		return
	}
}

func FakeLoginAndAuth(request *http.Request) (*User, error) {
	user, err := NewUser("fake_user_login", "12345", "mail@mail.ru", "yasher")
	if err != nil {
		return nil, err
	}
	session := NewSession()
	session.user = user
	err = session.Save()
	request.AddCookie(&http.Cookie{
		Name:   "sid",
		Secure: true,
		Value:  session.sid})
	return user, err
}

func TestGetProfile(t *testing.T) {
	InitModels()
	request, err := http.NewRequest("GET", "http://localhost/api/profile", nil)
	expectedBody := `{"type":"usinfo","status":"success","payload":{"login":"fake_user_login","email":"mail@mail.ru","name":"yasher"}}`

	response := httptest.NewRecorder()
	_, err = FakeLoginAndAuth(request)

	if err != nil {
		t.Fatal(err.Error())
	}

	router := NewRouter()
	router.ServeHTTP(response, request)
	result, _ := ioutil.ReadAll(response.Body)

	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result)
	}

}

func TestUpdateProfile(t *testing.T) {
	InitModels()
	body := strings.NewReader(`{
		"password" : "qweqwe234234&62342=",
		"name": "new name" }`)
	request, err := http.NewRequest("PUT", "http://localhost/api/profile", body)
	expectedBody := `{"type":"usinfo","status":"success","payload":{"login":"fake_user_login","email":"mail@mail.ru","name":"new name"}}`

	response := httptest.NewRecorder()
	user, err := FakeLoginAndAuth(request)

	if err != nil {
		t.Fatal(err.Error())
	}

	router := NewRouter()
	router.ServeHTTP(response, request)
	result, _ := ioutil.ReadAll(response.Body)

	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result)
	}
	user, _ = GetUserByLogin("fake_user_login")
	if user.name != "new name" {
		t.Errorf("Wrong name\n Expected:new name\nGot:%s", user.name)
	}
	if user.passwordHash != "qweqwe234234&62342=" {
		t.Errorf("Wrong passeord hash\n Expected:qweqwe234234&62342=\nGot:%s", user.passwordHash)
	}
}

func TestGetLeaderboard(t *testing.T) {
	InitModels()

	for i := 1; i <= 27; i++ {
		NewUser("npc_"+string(i), "12345", "mail"+string(i)+"@mail.ru", "Nick #"+string(i))
	}
	request, err := http.NewRequest("GET", "http://localhost/api/leaderboard/1", nil)
	expectedBody := `NOT IMPLEMENTED!` //TODO: Implement this

	response := httptest.NewRecorder()
	_, err = FakeLoginAndAuth(request)

	if err != nil {
		t.Fatal(err.Error())
	}

	router := NewRouter()
	router.ServeHTTP(response, request)
	result, _ := ioutil.ReadAll(response.Body)

	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result)
	}

}

func TestGetLeaderboardTooBigPage(t *testing.T) {
	InitModels()

	for i := 1; i <= 27; i++ {
		NewUser("npc_"+string(i), "12345", "mail"+string(i)+"@mail.ru", "Nick #"+string(i))
	}
	request, err := http.NewRequest("GET", "http://localhost/api/leaderboard/31", nil)
	expectedBody := `{"type":"uslist","status":"error","payload":{"message":"Not enough users"}}`

	response := httptest.NewRecorder()
	_, err = FakeLoginAndAuth(request)

	if err != nil {
		t.Fatal(err.Error())
	}

	router := NewRouter()
	router.ServeHTTP(response, request)
	result, _ := ioutil.ReadAll(response.Body)

	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result)
	}

}
