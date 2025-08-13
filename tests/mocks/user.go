package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/go-microservice-lib/pagination"
	"github.com/javiorfo/go-microservice-lib/response/backend"
	"github.com/javiorfo/nilo"
	"github.com/stretchr/testify/mock"
)

// Mock Repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByCode(ctx context.Context, code uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, code)
	if user, ok := args.Get(0).(*model.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (nilo.Optional[model.User], error) {
	args := m.Called(ctx, username)
	if user, ok := args.Get(0).(nilo.Optional[model.User]); ok {
		return user, args.Error(1)
	}
	return nilo.Empty[model.User](), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, queryFilter pagination.QueryFilter) ([]model.User, error) {
	args := m.Called(ctx, queryFilter)
	if users, ok := args.Get(0).([]model.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context, queryFilter pagination.QueryFilter) (int64, error) {
	args := m.Called(ctx, queryFilter)
	if users, ok := args.Get(0).(int64); ok {
		return users, args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Mock Service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FindAll(ctx context.Context, queryFilter pagination.QueryFilter) ([]model.User, error) {
	args := m.Called(ctx, queryFilter)
	if users, ok := args.Get(0).([]model.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) Count(ctx context.Context, queryFilter pagination.QueryFilter) (int64, error) {
	args := m.Called(ctx, queryFilter)
	if users, ok := args.Get(0).(int64); ok {
		return users, args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *MockUserService) Create(ctx context.Context, user *model.User, permission string) backend.Error {
	args := m.Called(ctx, user, permission)
	if be, ok := args.Get(0).(backend.Error); ok {
		return be
	}
	return nil
}

func (m *MockUserService) Login(ctx context.Context, username, password string) (string, backend.Error) {
	args := m.Called(ctx, username, password)
	if users, ok := args.Get(0).(string); ok {
		return users, nil
	}
	return "", args.Get(0).(backend.Error)
}
