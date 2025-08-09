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

func (userDb *UserDB) From(user model.User) {
	userDb.ID = user.ID
	userDb.Username = user.Username
	userDb.Email = user.Email
	userDb.Status = user.Status
	userDb.Permission.From(user.Permission)
}

func (userDb UserDB) Into() model.User {
	return model.User{
		ID:       userDb.ID,
		Code:     userDb.Code,
		Username: userDb.Username,
		Email:    userDb.Email,
		Status:   userDb.Status,
	}
}

/* func (u UserEntity) VerifyPassword(password string) bool {
	hashedInputPassword := pwd.Hash(password, u.Salt)
	return hashedInputPassword == u.Password
} */
