package entities

import "github.com/javiorfo/fiber-micro/application/domain/model"

type RoleDB struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"unique"`
	// Permissions []PermissionDB `json:"permissions" gorm:"many2many:permissions_roles;"`
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
