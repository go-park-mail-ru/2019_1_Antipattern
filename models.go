package main

import (
	"errors"
	"sort"

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
var uuidUserIndex map[uint32]string
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
	login, exists := uuidUserIndex[uuid]
	if !exists {
		return nil, errors.New("wrong uuid")
	}

	user, ok := users[login]
	if !ok {
		return nil, errors.New("uuid-login match error")
	}

	return &user, nil
}

func GetUserByLogin(login string) (*User, error) {
	user, exists := users[login]
	if !exists {
		return nil, errors.New("wrong login")
	}

	return &user, nil
}

func GetSession(id string) (*Session, error) {
	session, exists := sessions[id]
	if !exists {
		return nil, errors.New("Wrong sid")
	}
	return &session, nil
}

func GetUsers(count, page int) ([]User, error) {
	if page < 1 {
		return nil, errors.New("invalid page number")
	}

	min := count*(page-1)
	if min >= len(users) {
		return nil, errors.New("not enough users")
	}

	//var max uint = uint(math.Max(float64(count*page), float64(len(users))))

	max := count * page
	if max > len(users) {
		max = len(users)
	}

	keySlice := make([]string, 0, len(users))
	for k := range users {
		keySlice = append(keySlice, k)
	}
	sort.Strings(keySlice)

	userSlice := make([]User, 0, len(users))
	for _, v := range keySlice {
		userSlice = append(userSlice, users[v])
	}

	return userSlice[min:max], nil
}

func (session *Session) Delete() error {
	delete(sessions, session.sid)
	return nil
}

func (user *User) Delete() error {
	delete(uuidUserIndex, user.uuid)
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
	if login == "" {
		return nil, errors.New("login")
	}

	if password == "" {
		return nil, errors.New("password")
	}

	if email == "" {
		return nil, errors.New("email")
	}

	if name == "" {
		return nil, errors.New("name")
	}

	if _, ok := users[login]; ok {
		return nil, errors.New("user already exists")
	}
	user := User{
		uuid:         uuid.New().ID(),
		login:        login,
		passwordHash: password,
		email:        email,
		name:         name,
	}

	users[login] = user
	uuidUserIndex[user.uuid] = user.login
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

func InitModels() {
	users = make(map[string]User)
	uuidUserIndex = make(map[uint32]string)
	sessions = make(map[string]Session)
}
