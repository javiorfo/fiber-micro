package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/javiorfo/go-microservice-lib/security"
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
	Salt           string     `json:"-"`
	CreatedBy      string     `json:"-"`
	LastModifiedBy *string    `json:"-"`
	CreateDate     time.Time  `json:"-"`
	LastModified   *time.Time `json:"-"`
}

func NewUser(username string, email string, permission Permission, password string, createdBy string) User {
	return User{
		Username:   username,
		Email:      email,
		Password:   password,
		CreatedBy:  createdBy,
		Permission: permission,
		Status:     UserStatusActive,
	}
}

func (u User) VerifyPassword(password string) bool {
	hashedInputPassword := security.Hash(password, u.Salt)
	return hashedInputPassword == u.Password
}

type UserStatus = string

const (
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusInactive UserStatus = "INACTIVE"
	UserStatusBlocked  UserStatus = "BLOCKED"
)

var ValidateStatus = validation.NewEnumValidator(
	validation.Tag("status"),
	validation.JsonField("status"),
	UserStatusActive, UserStatusInactive, UserStatusBlocked)
