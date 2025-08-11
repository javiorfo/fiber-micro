package model

type Role struct {
	ID   uint   `json:"-"`
	Name string `json:"name"`
}

func NewRole(name string) *Role {
	return &Role{Name: name}
}
