package service

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/adapter/database/entities"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/fiber-micro/application/domain/service"
	"github.com/javiorfo/fiber-micro/tests/mocks"
	"github.com/javiorfo/go-microservice-lib/pagination"
	"github.com/javiorfo/go-microservice-lib/response"
	"github.com/javiorfo/go-microservice-lib/security"
	"github.com/javiorfo/nilo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_DURATION", "300")
	os.Setenv("JWT_SECRET_KEY", uuid.NewString())
	os.Setenv("JWT_ISSUER", "mock")
	os.Setenv("JWT_AUDIENCE", "mock")

	exitCode := m.Run()

	os.Unsetenv("JWT_DURATION")
	os.Unsetenv("JWT_SECRET_KEY")
	os.Unsetenv("JWT_ISSUER")
	os.Unsetenv("JWT_AUDIENCE")

	os.Exit(exitCode)
}

var userRepo = new(mocks.MockUserRepository)
var permRepo = new(mocks.MockPermissionRepository)
var userService = service.NewUserService(userRepo, permRepo)

func TestUserCreate(t *testing.T) {
	ctx := context.Background()
	perm := model.Permission{
		ID:    1,
		Name:  "PERM",
		Roles: []model.Role{{ID: 1, Name: "ROLE_1"}},
	}
	permRepo.On("FindByName", ctx, mock.Anything).Return(nilo.Some(perm), nil)

	user := model.User{
		ID:        1,
		Code:      uuid.New(),
		Username:  "Javi",
		Email:     "mail",
		Password:  "1234",
		Status:    "ACTIVE",
		CreatedBy: "auditor",
	}

	userRepo.On("Create", ctx, &user).Return(nil)

	backErr := userService.Create(ctx, &user, perm.Name)

	assert.Nil(t, backErr)
	userRepo.AssertExpectations(t)
	permRepo.AssertExpectations(t)
}

func TestUserLogin(t *testing.T) {
	username := "javi"
	password := "1234"
	salt, _ := security.GenerateSalt()
	hashed := security.Hash(password, salt)

	t.Logf("Hashed password %s, salt %s", hashed, salt)

	user := model.User{
		Code:     uuid.New(),
		Username: username,
		Password: hashed,
		Salt:     salt,
	}

	ctx := context.Background()
	userRepo.On("FindByUsername", ctx, mock.Anything).Return(nilo.Some(user), nil)

	token, backendErr := userService.Login(ctx, username, password)

	assert.Nil(t, backendErr)
	assert.NotEmpty(t, token)

	_, backendErr = userService.Login(ctx, username, "12345")

	resp, ok := backendErr.(*response.ResponseError)

	assert.True(t, ok)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Get().Code, response.ErrorCode("FIBER-MICRO-003"))
	userRepo.AssertExpectations(t)
}

func TestUserFindAll(t *testing.T) {
	queryFilter := entities.NewUserFilter(pagination.DefaultPage(), "Javi", "", "")
	expected := []model.User{
		{ID: 1},
		{ID: 2},
	}

	ctx := context.Background()
	userRepo.On("FindAll", ctx, queryFilter).Return(expected, nil)

	result, err := userService.FindAll(ctx, queryFilter)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	userRepo.AssertExpectations(t)
}
