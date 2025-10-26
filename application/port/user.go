package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/go-microservice-lib/response/backend"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/nilo"
)

type UserRepository interface {
	FindByCode(ctx context.Context, code uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (nilo.Option[model.User], error)
	FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[model.User], error)
	Create(ctx context.Context, user *model.User) error
}

type UserService interface {
	FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[model.User], error)
	Create(ctx context.Context, user *model.User, permission string) backend.Error
	Login(ctx context.Context, username, password string) (string, backend.Error)
}
