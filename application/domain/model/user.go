package model

import (
	"github.com/google/uuid"
	"github.com/javiorfo/go-microservice-lib/validation"
)

type User struct {
	ID         uint       `json:"-"`
	Code       uuid.UUID  `json:"code"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Permission Permission `json:"permission"`
	Status     Status     `json:"status"`
	Password   string     `json:"-"`
}

type Status = string

const (
	active   Status = "ACTIVE"
	inactive Status = "INACTIVE"
	blocked  Status = "BLOCKED"
)

var ValidateStatus = validation.NewEnumValidator("status", "status", active, inactive, blocked)
