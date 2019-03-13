package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"./models"
)

// User register

func CheckSessionSetCookie(t *testing.T, user models.User, w *httptest.ResponseRecorder) {
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

	session, err := models.GetSession(sessionID.Value)
	if err != nil {
		t.Errorf("Can't get session!\n%s", err.Error())
		return
	}
	if session.User.Uuid != user.Uuid {
		t.Errorf("Session uuid is wrong.\nExpected:%s\nGot:%s", string(user.Uuid), string(session.User.Uuid))
		return
	}
}
func TestRegister(t *testing.T) {
	models.InitModels()

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
	expectedBody := `{"type":"reg","status":"success","payload":{"login":"user_login","email":"death.pa_cito@mail.yandex.ru","name":"Gamer #23 @790-_%","score":20}}`

	w := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(w, r)

	result, _ := ioutil.ReadAll(w.Body)
	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result)
		return
	}
	newUser, _ := models.GetUserByLogin("user_login")

	if newUser.Login != "user_login" || newUser.PasswordHash != "qweqwe234234&62342=" ||
		newUser.Email != "death.pa_cito@mail.yandex.ru" || newUser.Name != "Gamer #23 @790-_%" {
		t.Errorf("Wrong user in db")
		return
	}
	CheckSessionSetCookie(t, *newUser, w)
}
func TestRegisterAlreadyRegistered(t *testing.T) {
	models.InitModels()
	_, err := models.NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru", "kek")
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
	expectedBody := `{"type":"reg","status":"error","payload":{"message":"user already exists","field":"login"}}`

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
	models.InitModels()
	user, err := models.NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru", "kek")
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
	expectedBody := `{"type":"log","status":"success","payload":{"login":"user_login","email":"death.pa_cito@mail.yandex.ru","name":"kek","score":20}}`

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
	models.InitModels()
	_, err := models.NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru", "kek")
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
	expectedBody := `{"type":"log","status":"error","payload":{"message":"incorrect password","field":"password"}}`

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
	models.InitModels()
	_, err := models.NewUser("user_login", "1235689", "death.pa_cito@mail.yandex.ru", "kek")
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
	expectedBody := `{"type":"log","status":"error","payload":{"message":"incorrect login","field":"login"}}`

	w := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(w, r)

	result, _ := ioutil.ReadAll(w.Body)
	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", expectedBody, result)
		return
	}
}

func FakeLoginAndAuth(request *http.Request) (*models.User, error) {
	user, err := models.NewUser("fake_user_login", "12345", "mail@mail.ru", "yasher")
	if err != nil {
		return nil, err
	}
	session := models.NewSession()
	session.User = user
	err = session.Save()
	request.AddCookie(&http.Cookie{
		Name:   "sid",
		Secure: true,
		Value:  session.Sid})
	return user, err
}

func TestGetProfile(t *testing.T) {
	models.InitModels()
	request, err := http.NewRequest("GET", "http://localhost/api/profile", nil)
	expectedBody := `{"type":"usinfo","status":"success","payload":{"login":"fake_user_login","email":"mail@mail.ru","name":"yasher","score":20}}`

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
	models.InitModels()
	body := strings.NewReader(`{
		"password" : "qweqwe234234&62342=",
		"name": "new name" }`)
	request, err := http.NewRequest("PUT", "http://localhost/api/profile", body)
	expectedBody := `{"type":"usinfo","status":"success","payload":{"login":"fake_user_login","email":"mail@mail.ru","name":"new name","score":20}}`

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
	user, _ = models.GetUserByLogin("fake_user_login")
	if user.Name != "new name" {
		t.Errorf("Wrong name\n Expected:new name\nGot:%s", user.Name)
	}
	if user.PasswordHash != "qweqwe234234&62342=" {
		t.Errorf("Wrong passeord hash\n Expected:qweqwe234234&62342=\nGot:%s", user.PasswordHash)
	}
}

func TestGetLeaderboard(t *testing.T) {
	models.InitModels()

	for i := 1; i <= 27; i++ {
		models.NewUser("npc_"+strconv.Itoa(i), "12345", "mail"+strconv.Itoa(i)+"@mail.ru", "Nick #"+strconv.Itoa(i))
	}
	request, err := http.NewRequest("GET", "http://localhost/api/leaderboard/1", nil)
	expectedBody := `{"type":"uslist","status":"success","payload":{"users":[{"name":"yasher","score":20},{"name":"Nick #1","score":20},{"name":"Nick #10","score":20},{"name":"Nick #11","score":20},{"name":"Nick #12","score":20},{"name":"Nick #13","score":20},{"name":"Nick #14","score":20},{"name":"Nick #15","score":20},{"name":"Nick #16","score":20},{"name":"Nick #17","score":20}],"count":28}}`

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
	models.InitModels()

	for i := 1; i <= 27; i++ {
		models.NewUser("npc_"+string(i), "12345", "mail"+string(i)+"@mail.ru", "Nick #"+string(i))
	}
	request, err := http.NewRequest("GET", "http://localhost/api/leaderboard/31", nil)
	expectedBody := `{"type":"uslist","status":"error","payload":{"message":"not enough users"}}`

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
