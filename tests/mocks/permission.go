package mocks

import (
	"context"

	"github.com/javiorfo/fiber-micro/application/domain/model"
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
