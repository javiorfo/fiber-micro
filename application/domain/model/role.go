package model

type Role struct {
	ID          uint         `json:"-"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}
