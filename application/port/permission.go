package port

import (
	"context"

	"github.com/javiorfo/fiber-micro/application/domain/model"
)

type PermissionRepository interface {
	Create(ctx context.Context, user *model.Permission) error
}
