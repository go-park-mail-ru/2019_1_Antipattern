package main

type Response struct {
	Type    string      `json:"type"`
	Status  string      `json:"status"`
	Payload interface{} `json:"payload,omitempty"`
}

type UserDataPayload struct {
	Login      string `json:"login,omitempty"`
	Email      string `json:"email,omitempty"`
	Name       string `json:"name,omitempty"`
	AvatarPath string `json:"avatar,omitempty"`
}

type ErrorPayload struct {
	Message string `json:"message,omitempty"`
	Field   string `json:"field,omitempty"`
}

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}
