package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/go-microservice-lib/pagination"
)

type UserRepository interface {
	FindByCode(ctx context.Context, code uuid.UUID) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindAll(ctx context.Context, page pagination.Page, info string) ([]model.User, error)
	Create(ctx context.Context, user *model.User, auditor string) error
}

type UserService interface {
	FindByCode(ctx context.Context, code uuid.UUID) (*model.User, error)
	FindAll(ctx context.Context, page pagination.Page, info string) ([]model.User, error)
	Create(user *model.User, auditor string) (*string, error)
	Login(username, password string) (string, error)
}
