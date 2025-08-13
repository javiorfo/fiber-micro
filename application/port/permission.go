package port

import (
	"context"

	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/nilo"
)

type PermissionRepository interface {
	Create(ctx context.Context, user *model.Permission) error
	FindByName(ctx context.Context, name string) (nilo.Optional[model.Permission], error)
}
