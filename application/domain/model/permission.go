package model

type Permission struct {
	ID    uint   `json:"-"`
	Name  string `json:"name"`
	Roles []Role `json:"roles"`
}

func NewPermission(name string, roles []Role) *Permission {
	return &Permission{Name: name, Roles: roles}
}
