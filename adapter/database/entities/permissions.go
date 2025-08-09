package entities

import "github.com/javiorfo/fiber-micro/application/domain/model"

type PermissionDB struct {
	ID    uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string   `json:"name" gorm:"unique"`
	Roles []RoleDB `json:"roles" gorm:"many2many:permissions_roles;"`
}

func (PermissionDB) TableName() string {
	return "permissions"
}

func (permDb *PermissionDB) From(perm model.Permission) {
	permDb.ID = perm.ID
	permDb.Name = perm.Name
}

func (permDb PermissionDB) Into() model.Permission {
	return model.Permission{
		Name: permDb.Name,
	}
}
