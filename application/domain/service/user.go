package service

import (
	"context"

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

	permission, err := service.permRepository.FindByName(ctx, permName)
	if err != nil {
		return errors.PermissionNotFound(span)
	}

	salt, err := security.GenerateSalt()
	if err != nil {
		return backend.GenericError(span, err)
	}

	hashedPassword := security.Hash(user.Password, salt)

	// Set user values
	user.Permission = *permission
	user.Password = hashedPassword
	user.Salt = salt

	return backend.GenericError(span, service.userRepository.Create(ctx, user))
}

func (service *userService) Login(ctx context.Context, username string, password string) (string, backend.Error) {
	_, span := service.tracer.Start(ctx, tracing.Name())
	defer span.End()

	user, err := service.userRepository.FindByUsername(ctx, username)
	if err != nil {
		return "", errors.UserNotFound(span)
	}

	if user.VerifyPassword(password) {
		roles := steams.OfSlice(user.Permission.Roles).MapToString(func(r model.Role) string {
			return r.Name
		}).Collect()

		tokenPermission := security.TokenPermission{Name: user.Permission.Name, Roles: roles}

		token, err := security.CreateToken(tokenPermission, username)
		if err != nil {
			return "", backend.GenericError(span, err)
		}
		return token, nil
	}

	return "", errors.CredentialsError(span)
}

func (service *userService) FindAll(ctx context.Context, filter pagination.QueryFilter) ([]model.User, error) {
	_, span := service.tracer.Start(ctx, tracing.Name())
	defer span.End()

	log.Infof("%sFilter %+v", tracing.Log(span), filter)

	return service.userRepository.FindAll(ctx, filter)
}
