package main

import (
	"errors"

	"github.com/google/uuid"
)

type Model interface {
	Delete() error
	Save() error
}

type User struct {
	uuid         uint32
	login        string
	passwordHash string
	email        string
	name         string
	avatar       string
}

type Session struct {
	sid  string
	user *User
}

var users map[string]User
var sessions map[string]Session

func (session *Session) Save() error {
	sessions[session.sid] = *session
	return nil
}

func (user *User) Save() error {
	users[user.login] = *user
	return nil
}

func GetUser(uuid uint32) (*User, error) {
	return nil, nil
}

func GetUsers(count, page int) ([]User, error) {
	return nil, nil
}

func (session *Session) Delete() error {
	delete(sessions, session.sid)
	return nil
}
func (user *User) Delete() error {
	delete(users, user.login)
	return nil
}

func NewSession() *Session {
	id := uuid.New().String()
	session := Session{
		sid:  id,
		user: nil,
	}

	sessions[id] = session
	return &session
}

func NewUser(login string, password string, email string, name string) (*User, error) {
	if _, ok := users[login]; ok {
		return nil, errors.New("User already exists " + login)
	}
	user := User{
		uuid:         uuid.New().ID(),
		login:        login,
		passwordHash: password,
		email:        email,
		name:         name,
	}

	users[login] = user
	return &user, nil
}

func Auth(login string, password string) (*User, error) {
	user, ok := users[login]
	if !ok {
		return nil, errors.New("login")
	}
	if user.passwordHash != password {
		return nil, errors.New("password")
	}

	return &user, nil
}
