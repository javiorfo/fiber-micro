package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/go-microservice-lib/pagination"
	"github.com/javiorfo/go-microservice-lib/response/backend"
	"github.com/javiorfo/nilo"
)

type UserRepository interface {
	FindByCode(ctx context.Context, code uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (nilo.Optional[model.User], error)
	FindAll(ctx context.Context, queryFilter pagination.QueryFilter) ([]model.User, error)
	Count(ctx context.Context, queryFilter pagination.QueryFilter) (int64, error)
	Create(ctx context.Context, user *model.User) error
}

type UserService interface {
	FindAll(ctx context.Context, queryFilter pagination.QueryFilter) ([]model.User, error)
	Count(ctx context.Context, queryFilter pagination.QueryFilter) (int64, error)
	Create(ctx context.Context, user *model.User, permission string) backend.Error
	Login(ctx context.Context, username, password string) (string, backend.Error)
}
