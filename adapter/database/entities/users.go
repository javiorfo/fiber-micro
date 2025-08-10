package entities

import (
	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/go-microservice-lib/auditory"
)

type UserDB struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Code         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Username     string    `gorm:"not null"`
	Email        string    `gorm:"not null"`
	PermissionID uint      `gorm:"not null"`
	Permission   PermissionDB
	Password     string `gorm:"not null"`
	Salt         string `gorm:"not null"`
	Status       string `gorm:"not null"`
	auditory.Auditable
}

func (UserDB) TableName() string {
	return "users"
}

func (userDB *UserDB) From(user model.User) {
	userDB.ID = user.ID
	userDB.Username = user.Username
	userDB.Email = user.Email
	userDB.Status = user.Status
	userDB.Permission.From(user.Permission)

	if user.LastModifiedBy == nil {
		userDB.Auditable = *auditory.New(user.CreatedBy)
	} else {
		userDB.Auditable.Update(user.LastModifiedBy)
	}
}

func (userDB UserDB) Into() *model.User {
	return &model.User{
		ID:             userDB.ID,
		Code:           userDB.Code,
		Username:       userDB.Username,
		Email:          userDB.Email,
		Status:         userDB.Status,
		CreatedBy:      userDB.CreatedBy,
		CreateDate:     userDB.CreateDate,
		LastModifiedBy: userDB.LastModifiedBy,
		LastModified:   userDB.LastModified,
	}
}

/* func (u UserEntity) VerifyPassword(password string) bool {
	hashedInputPassword := pwd.Hash(password, u.Salt)
	return hashedInputPassword == u.Password
} */
