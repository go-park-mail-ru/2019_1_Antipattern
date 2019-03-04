package main

type RegResponce struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Field   string `json:"field,omitempty"`
}

type LogResponce struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Field   string `json:"field,omitempty"`
}

type AuthResponce struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type UserInfoResponce struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Payload struct {
		Login      string
		Email      string
		Name       string
		AvatarPath string
	} `json:"payload,omitempty"`
}
