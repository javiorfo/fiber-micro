package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/go-microservice-lib/response/backend"
	"github.com/javiorfo/gormen/pagination"
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

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (nilo.Option[model.User], error) {
	args := m.Called(ctx, username)
	if user, ok := args.Get(0).(nilo.Option[model.User]); ok {
		return user, args.Error(1)
	}
	return nilo.None[model.User](), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[model.User], error) {
	args := m.Called(ctx, pageable)
	if users, ok := args.Get(0).(*pagination.Page[model.User]); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Mock Service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[model.User], error) {
	args := m.Called(ctx, pageable)
	if users, ok := args.Get(0).(*pagination.Page[model.User]); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
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
	var token string
	if args.Get(0) != nil {
		token = args.Get(0).(string)
	}

	var err backend.Error
	if args.Get(1) != nil {
		err = args.Get(1).(backend.Error)
	}

	return token, err
}
