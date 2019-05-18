package user_data

type User struct {
	Uid    string
	Login  string
	Avatar string
}

// Provider is an interface for api providers
type Provider interface {
	GetUsers(ids []string) ([]*User, error)
}
