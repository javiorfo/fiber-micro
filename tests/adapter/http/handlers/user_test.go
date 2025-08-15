package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/javiorfo/fiber-micro/adapter/database/entities"
	"github.com/javiorfo/fiber-micro/adapter/http/request"
	srvResp "github.com/javiorfo/fiber-micro/adapter/http/response"
	"github.com/javiorfo/fiber-micro/adapter/http/routes"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	be "github.com/javiorfo/fiber-micro/application/domain/service/errors"
	"github.com/javiorfo/fiber-micro/tests/mocks"
	"github.com/javiorfo/go-microservice-lib/pagination"
	"github.com/javiorfo/go-microservice-lib/response"
	"github.com/javiorfo/go-microservice-lib/response/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
)

var app *fiber.App
var mockSecurity *mocks.MockAuthorizer
var mockUserService *mocks.MockUserService

func TestMain(m *testing.M) {
	app = fiber.New()
	mockSecurity = new(mocks.MockAuthorizer)
	mockUserService = new(mocks.MockUserService)

	routes.User(app, mockSecurity, mockUserService)

	os.Exit(m.Run())
}

func TestLogin(t *testing.T) {
	tracer := otel.Tracer("Login")
	ctx, span := tracer.Start(context.Background(), "mockpath")
	defer span.End()

	loginReq := request.LoginRequest{
		Username: "javi",
		Password: "12334",
	}

	t.Run("Successful", func(t *testing.T) {
		mockUserService.ExpectedCalls = nil
		mockUserService.On("Login", ctx, mock.Anything, mock.Anything).Return("token", nil)

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(fiber.MethodPost, "/users/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseBody srvResp.LoginResponse
		json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.Equal(t, "javi", responseBody.Username)
		assert.Equal(t, "token", responseBody.Token)

		mockUserService.AssertExpectations(t)
	})
	
	t.Run("Credentials error", func(t *testing.T) {
		mockUserService.ExpectedCalls = nil
		mockUserService.On("Login", ctx, mock.Anything, mock.Anything).Return("", be.CredentialsError(span))

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(fiber.MethodPost, "/users/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var responseBody response.ResponseError
		json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.Equal(t, "FIBER-MICRO-003", responseBody.Get().Code)
		assert.Equal(t, "Username or password incorrect!", responseBody.Get().Message)

		mockUserService.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	tracer := otel.Tracer("CreateUser")
	ctx, span := tracer.Start(context.Background(), "mockpath")
	defer span.End()

	userRequest := request.CreateUserRequest{
		Username:   "javi",
		Status:     "ACTIVE",
		Password:   "1234",
		Permission: "PERM",
		Email:      "some@mail",
	}

	t.Run("Successful", func(t *testing.T) {
		mockUserService.ExpectedCalls = nil
		mockUserService.On("Create", ctx, mock.Anything, "PERM").Return(nil)

		body, _ := json.Marshal(userRequest)
		req := httptest.NewRequest(fiber.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var responseBody srvResp.CreateUserResponse
		json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.Equal(t, "javi", responseBody.User.Username)

		mockUserService.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mockUserService.ExpectedCalls = nil
		body := `{ "invalid": 10 }`
		req := httptest.NewRequest(fiber.MethodPost, "/users", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockUserService.ExpectedCalls = nil
		mockUserService.On("Create", ctx, mock.Anything, "PERM").Return(backend.InternalError(span, errors.New("service error")))

		body, _ := json.Marshal(userRequest)
		req := httptest.NewRequest(fiber.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUserService.AssertExpectations(t)
	})
}

func TestHandlerFindAll(t *testing.T) {
	tracer := otel.Tracer("FindAllUsers")
	ctx, span := tracer.Start(context.Background(), "mockpath")
	defer span.End()

	t.Run("Successful", func(t *testing.T) {
		page := pagination.Page{Page: 1, Size: 10, SortBy: "username", SortOrder: "asc"}
		userFilter := entities.NewUserFilter(page, "", "", "")
		mockUserService.On("FindAll", ctx, userFilter).Return([]model.User{{ID: 1, Username: "javi"}}, nil)
		mockUserService.On("Count", ctx, userFilter).Return(int64(1), nil)

		req := httptest.NewRequest("GET", "/users?page=1&size=10&sortBy=username&sortOrder=asc", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseBody response.PaginationResponse[model.User]
		json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.Equal(t, int64(1), responseBody.Pagination.Total)
		assert.Equal(t, "javi", responseBody.Elements[0].Username)

		mockUserService.AssertExpectations(t)
	})

	t.Run("DB Error", func(t *testing.T) {
		mockUserService.On("FindAll", ctx, mock.Anything).Return(nil, errors.New("data source error"))

		req := httptest.NewRequest("GET", "/users?page=1&size=10&sortBy=id&sortOrder=asc", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		mockUserService.AssertExpectations(t)
	})

	t.Run("Pagination Bad Request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users?page=invalid&size=10&sortBy=id&sortOrder=asc", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
