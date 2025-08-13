package request

import (
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/go-microservice-lib/validation"
)

type CreateUserRequest struct {
	Username   string `json:"username" validate:"required,notblank"`
	Email      string `json:"email" validate:"required,notblank"`
	Password   string `json:"password" validate:"required,notblank"`
	Status     string `json:"status" validate:"required,status"`
	Permission string `json:"permission" validate:"required,notblank"`
}

func (cur CreateUserRequest) Into(auditor string) model.User {
	return model.User{
		Username:  cur.Username,
		Email:     cur.Email,
		Password:  cur.Password,
		Status:    cur.Status,
		CreatedBy: auditor,
	}
}

var ValidateUserNotBlank = validation.NewNotBlankValidator("(username, email, password, status and permission)")

type LoginRequest struct {
	Username string `json:"username" validate:"required,notblank"`
	Password string `json:"password" validate:"required,notblank"`
}
