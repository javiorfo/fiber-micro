package entities

import (
	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/go-microservice-lib/auditory"
	"github.com/javiorfo/go-microservice-lib/pagination"
	"gorm.io/gorm"
)

type UserDB struct {
	ID           uint         `gorm:"primaryKey;autoIncrement"`
	Code         uuid.UUID    `gorm:"not null"`
	Username     string       `gorm:"not null"`
	Email        string       `gorm:"not null"`
	PermissionID uint         `gorm:"not null"`
	Permission   PermissionDB `gorm:"column:permission_id;not null"`
	Password     string       `gorm:"not null"`
	Salt         string       `gorm:"not null"`
	Status       string       `gorm:"not null"`
	auditory.Auditable
}

func (UserDB) TableName() string {
	return "users"
}

func (userDB *UserDB) BeforeCreate(tx *gorm.DB) (err error) {
	userDB.Code = uuid.New()
	return
}

func (userDB *UserDB) From(user model.User) {
	userDB.ID = user.ID
	userDB.Username = user.Username
	userDB.Email = user.Email
	userDB.Status = user.Status
	userDB.Password = user.Password
	userDB.Salt = user.Salt
	userDB.Permission.From(user.Permission)

	if user.LastModifiedBy == nil {
		userDB.Auditable = auditory.New(user.CreatedBy)
	} else {
		userDB.Auditable.Update(user.LastModifiedBy)
	}
}

func (userDB UserDB) Into() model.User {
	return model.User{
		ID:             userDB.ID,
		Code:           userDB.Code,
		Username:       userDB.Username,
		Email:          userDB.Email,
		Status:         userDB.Status,
		Password:       userDB.Password,
		Salt:           userDB.Salt,
		CreatedBy:      userDB.CreatedBy,
		CreateDate:     userDB.CreateDate,
		LastModifiedBy: userDB.LastModifiedBy,
		LastModified:   userDB.LastModifiedDate,
		Permission:     userDB.Permission.Into(),
	}
}

type userFilter struct {
	Username       string `filter:"users.username = ?"`
	PermissionName string `filter:"permissions.name = ?"`
	CreateDate     string `filter:"create_date > ?;type:time.Time"`
	pagination.Page
}

func (uf userFilter) PaginateAndFilter(db *gorm.DB) (*gorm.DB, error) {
	return pagination.PaginateAndFilter(db, uf.Page, uf)
}

func (uf userFilter) FilterOnly(db *gorm.DB) (*gorm.DB, error) {
	return pagination.FilterOnly(db, uf)
}

func NewUserFilter(page pagination.Page, username, permissionName, createDate string) *userFilter {
	return &userFilter{
		Username:       username,
		PermissionName: permissionName,
		CreateDate:     createDate,
		Page:           page,
	}
}
