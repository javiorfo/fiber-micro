package model

type Permission struct {
	ID    uint   `json:"-"`
	Name  string `json:"name"`
	Roles []Role `json:"roles"`
}
