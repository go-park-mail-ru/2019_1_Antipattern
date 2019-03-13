package models

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
	Uuid         uint32
	Login        string
	PasswordHash string
	Email        string
	Name         string
	Avatar       string
	Score        int
}

type Session struct {
	Sid  string
	User *User
}

var Users map[string]User
var UuidUserIndex map[uint32]string
var Sessions map[string]Session

func (session *Session) Save() error {
	Sessions[session.Sid] = *session
	return nil
}

func (user *User) Save() error {
	Users[user.Login] = *user
	return nil
}

func GetUser(uuid uint32) (*User, error) {
	login, exists := UuidUserIndex[uuid]
	if !exists {
		return nil, errors.New("wrong uuid")
	}

	user, ok := Users[login]
	if !ok {
		return nil, errors.New("uuid-login match error")
	}

	return &user, nil
}

func GetUserByLogin(login string) (*User, error) {
	user, exists := Users[login]
	if !exists {
		return nil, errors.New("wrong login")
	}

	return &user, nil
}

func GetSession(id string) (*Session, error) {
	session, exists := Sessions[id]
	if !exists {
		return nil, errors.New("Wrong sid")
	}
	return &session, nil
}

func GetUsers(count, page int) ([]User, error) {
	if page < 1 {
		return nil, errors.New("invalid page number")
	}

	min := count * (page - 1)
	if min >= len(Users) {
		return nil, errors.New("not enough users")
	}

	//var max uint = uint(math.Max(float64(count*page), float64(len(users))))

	max := count * page
	if max > len(Users) {
		max = len(Users)
	}

	keySlice := make([]string, 0, len(Users))
	for k := range Users {
		keySlice = append(keySlice, k)
	}
	sort.Strings(keySlice)

	userSlice := make([]User, 0, len(Users))
	for _, v := range keySlice {
		userSlice = append(userSlice, Users[v])
	}

	return userSlice[min:max], nil
}

func (session *Session) Delete() error {
	delete(Sessions, session.Sid)
	return nil
}

func (user *User) Delete() error {
	delete(UuidUserIndex, user.Uuid)
	delete(Users, user.Login)
	return nil
}

func NewSession() *Session {
	id := uuid.New().String()
	session := Session{
		Sid:  id,
		User: nil,
	}
	Sessions[id] = session
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

	if _, ok := Users[login]; ok {
		return nil, errors.New("user already exists")
	}
	user := User{
		Uuid:         uuid.New().ID(),
		Login:        login,
		PasswordHash: password,
		Email:        email,
		Name:         name,
		Score:        20,
	}

	Users[login] = user
	UuidUserIndex[user.Uuid] = user.Login
	return &user, nil
}

func Auth(login string, password string) (*User, error) {
	user, ok := Users[login]
	if !ok {
		return nil, errors.New("login")
	}
	if user.PasswordHash != password {
		return nil, errors.New("password")
	}

	return &user, nil
}

func GetUserCount() (int, error) {
	return len(Users), nil
}
func InitModels() {
	Users = make(map[string]User)
	UuidUserIndex = make(map[uint32]string)
	Sessions = make(map[string]Session)
}
