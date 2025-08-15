package service

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/fiber-micro/application/domain/service/errors"
	"github.com/javiorfo/fiber-micro/application/port"
	"github.com/javiorfo/go-microservice-lib/pagination"
	"github.com/javiorfo/go-microservice-lib/response/backend"
	"github.com/javiorfo/go-microservice-lib/security"
	"github.com/javiorfo/go-microservice-lib/tracing"
	"github.com/javiorfo/steams"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type userService struct {
	userRepository port.UserRepository
	permRepository port.PermissionRepository
	tracer         trace.Tracer
}

func NewUserService(ur port.UserRepository, pr port.PermissionRepository) port.UserService {
	return &userService{
		userRepository: ur,
		permRepository: pr,
		tracer:         otel.Tracer(tracing.Name()),
	}
}

func (service *userService) Create(ctx context.Context, user *model.User, permName string) backend.Error {
	_, span := service.tracer.Start(ctx, tracing.Name())
	defer span.End()

	permissionOpt, err := service.permRepository.FindByName(ctx, permName)
	if err != nil {
		return backend.InternalError(span, err)
	}

	if permissionOpt.IsPresent() {
		user.Permission = permissionOpt.Get()
	} else {
		return errors.PermissionNotFound(span)
	}

	salt, err := security.GenerateSalt()
	if err != nil {
		return backend.InternalError(span, err)
	}

	hashedPassword := security.Hash(user.Password, salt)
	user.Password = hashedPassword
	user.Salt = salt

	log.Info(tracing.LogInfo(span, "Hashed password created!"))

	err = service.userRepository.Create(ctx, user)
	if err != nil {
		return backend.InternalError(span, err)
	}

	return nil
}

func (service *userService) Login(ctx context.Context, username string, password string) (string, backend.Error) {
	_, span := service.tracer.Start(ctx, tracing.Name())
	defer span.End()

	userOpt, err := service.userRepository.FindByUsername(ctx, username)
	if err != nil {
		return "", backend.InternalError(span, err)
	}

	if userOpt.IsEmpty() {
		return "", errors.UserNotFound(span)
	}

	user := userOpt.Get()

	if user.VerifyPassword(password) {
		roles := steams.OfSlice(user.Permission.Roles).MapToString(func(r model.Role) string {
			return r.Name
		}).Collect()

		tokenPermission := security.TokenPermission{Name: user.Permission.Name, Roles: roles}

		token, err := security.CreateToken(tokenPermission, username)
		if err != nil {
			return "", backend.InternalError(span, err)
		}
		return token, nil
	}

	return "", errors.CredentialsError(span)
}

func (service *userService) FindAll(ctx context.Context, filter pagination.QueryFilter) ([]model.User, error) {
	_, span := service.tracer.Start(ctx, tracing.Name())
	defer span.End()

	log.Info(tracing.LogInfo(span, fmt.Sprintf("Filter %+v", filter)))

	return service.userRepository.FindAll(ctx, filter)
}

func (service *userService) Count(ctx context.Context, filter pagination.QueryFilter) (int64, error) {
	_, span := service.tracer.Start(ctx, tracing.Name())
	defer span.End()

	return service.userRepository.Count(ctx, filter)
}
