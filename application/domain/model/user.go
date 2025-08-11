package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/javiorfo/go-microservice-lib/validation"
)

type User struct {
	ID             uint       `json:"-"`
	Code           uuid.UUID  `json:"code"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	Permission     Permission `json:"permission"`
	Status         UserStatus `json:"status"`
	Password       string     `json:"-"`
	CreatedBy      string     `json:"-"`
	LastModifiedBy *string    `json:"-"`
	CreateDate     time.Time  `json:"-"`
	LastModified   *time.Time `json:"-"`
}

func NewUser(username string, email string, permission Permission, password string, createdBy string) User {
	return User{
		Username:   username,
		Email:      email,
		Permission: permission,
		Password:   password,
		CreatedBy:  createdBy,
	}
}

type UserStatus = string

const (
	active   UserStatus = "ACTIVE"
	inactive UserStatus = "INACTIVE"
	blocked  UserStatus = "BLOCKED"
)

var ValidateStatus = validation.NewEnumValidator(
	validation.Tag("status"),
	validation.JsonField("status"),
	active, inactive, blocked)
