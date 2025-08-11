package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/go-microservice-lib/pagination"
	"github.com/javiorfo/go-microservice-lib/response/backend"
)

type UserRepository interface {
	FindByCode(ctx context.Context, code uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindAll(ctx context.Context, filter pagination.QueryFilter) ([]model.User, error)
	Create(ctx context.Context, user *model.User) error
}

type UserService interface {
	FindAll(ctx context.Context, filter pagination.QueryFilter) ([]model.User, error)
	Create(ctx context.Context, user *model.User, permission string) backend.Error
	Login(ctx context.Context, username, password string) (string, backend.Error)
}
