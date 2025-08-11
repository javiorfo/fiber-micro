package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/go-microservice-lib/pagination"
)

type UserRepository interface {
	FindByCode(ctx context.Context, code uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindAll(ctx context.Context, filter pagination.GormFilter) ([]model.User, error)
	Create(ctx context.Context, user *model.User) error
}

type UserService interface {
	FindByCode(ctx context.Context, code uuid.UUID) (*model.User, error)
	FindAll(ctx context.Context, filter pagination.GormFilter) ([]model.User, error)
	Create(user *model.User) (*string, error)
	Login(username, password string) (string, error)
}
