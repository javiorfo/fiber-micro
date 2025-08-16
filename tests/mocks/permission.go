package mocks

import (
	"context"

	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/nilo"
	"github.com/stretchr/testify/mock"
)

// Mock Repository
type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) Create(ctx context.Context, perm *model.Permission) error {
	args := m.Called(ctx, perm)
	return args.Error(0)
}

func (m *MockPermissionRepository) FindByName(ctx context.Context, name string) (nilo.Option[model.Permission], error) {
	args := m.Called(ctx, name)
	if perm, ok := args.Get(0).(nilo.Option[model.Permission]); ok {
		return perm, args.Error(1)
	}
	return nilo.None[model.Permission](), args.Error(1)
}
