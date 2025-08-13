package response

import "github.com/javiorfo/fiber-micro/application/domain/model"

type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	Message  string `json:"message"`
}

func NewLoginResponse(username, token string) LoginResponse {
	return LoginResponse{username, token, "Access granted"}
}

type CreateUserResponse struct {
	User model.User `json:"user"`
}
