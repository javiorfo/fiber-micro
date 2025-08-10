package entities

import (
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/steams"
)

type PermissionDB struct {
	ID    uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string   `json:"name" gorm:"unique"`
	Roles []RoleDB `json:"roles" gorm:"many2many:permissions_roles;"`
}

func (PermissionDB) TableName() string {
	return "permissions"
}

func (permDB *PermissionDB) From(perm model.Permission) {
	permDB.ID = perm.ID
	permDB.Name = perm.Name
	permDB.Roles = steams.Mapping(steams.OfSlice(perm.Roles), func(role model.Role) RoleDB {
		var roleDB RoleDB
		roleDB.From(role)
		return roleDB
	}).Collect()
}

func (permDB PermissionDB) Into() model.Permission {
	return model.Permission{
		ID:   permDB.ID,
		Name: permDB.Name,
		Roles: steams.Mapping(steams.OfSlice(permDB.Roles), func(roleDB RoleDB) model.Role {
			return roleDB.Into()
		}).Collect(),
	}
}
