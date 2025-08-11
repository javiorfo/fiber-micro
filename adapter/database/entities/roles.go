package entities

import "github.com/javiorfo/fiber-micro/application/domain/model"

type RoleDB struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"unique"`
}

func (RoleDB) TableName() string {
	return "roles"
}

func (roleDB *RoleDB) From(role model.Role) {
	roleDB.ID = role.ID
	roleDB.Name = role.Name
}

func (roleDB RoleDB) Into() model.Role {
	return model.Role{
		ID:   roleDB.ID,
		Name: roleDB.Name,
	}
}
