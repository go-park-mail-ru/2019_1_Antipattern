package main

type Response struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Field   string `json:"field,omitempty"`
	Payload	ResponsePayload `json:"payload,omitempty"`
}

type ResponsePayload struct {
	Login      string `json:"login,omitempty"`
	Email      string `json:"email,omitempty"`
	Name       string `json:"name,omitempty"`
	AvatarPath string `json:"avatar,omitempty"`
}