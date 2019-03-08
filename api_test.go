package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// User register

func TestRegister(t *testing.T) {
	InitModels()
	t.Parallel()

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
	}
	expectedBody := `{"type":"reg","status":"success","payload":{"login":"user_login","email":"death.pa_cito@mail.yandex.ru","name":"Gamer #23 @790-_%"}}`
	w := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(w, r)

	result, _ := ioutil.ReadAll(w.Body)
	if strings.TrimSpace(string(result)) != expectedBody {
		t.Errorf("Wrong result\n Expected:%s\nGot:%s", result, expectedBody)
	}
	fmt.Println(w.Body)
}
