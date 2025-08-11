package port

import (
	"context"

	"github.com/javiorfo/fiber-micro/application/domain/model"
)

type RoleRepository interface {
	Create(ctx context.Context, user *model.Role) error
}
