package entities

type RoleDB struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"unique"`
	Permissions []PermissionDB `json:"permissions" gorm:"many2many:permissions_roles;"`
}

func (RoleDB) TableName() string {
	return "roles"
}
