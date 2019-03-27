package json_structs

type Response struct {
	Type    string      `json:"type"`
	Status  string      `json:"status"`
	Payload interface{} `json:"payload,omitempty"`
}

type UsersPayload struct {
	Users []UserDataPayload `json:"users"`
	Count int               `json:"count"`
}

type UserDataPayload struct {
	Login      string `json:"login,omitempty"`
	Email      string `json:"email,omitempty"`
	AvatarPath string `json:"avatar,omitempty"`
	Score      int    `json:"score"`
}

type ErrorPayload struct {
	Message string `json:"message,omitempty"`
	Field   string `json:"field,omitempty"`
}

type UsrRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LeaderboardRequest struct {
	Count int `json:"count"`
	Page  int `json:"page"`
}