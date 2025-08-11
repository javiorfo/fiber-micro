package mocks

import (
	"context"

	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/stretchr/testify/mock"
)

// Mock Repository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(ctx context.Context, role *model.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}
