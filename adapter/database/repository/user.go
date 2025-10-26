package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/adapter/database/entities"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/fiber-micro/application/port"
	"github.com/javiorfo/go-microservice-lib/tracing"
	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/converter"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/nilo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type userRepository struct {
	*gorm.DB
	tracer trace.Tracer
	gr     gormen.ReadRepository[model.User]
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &userRepository{
		DB:     db,
		tracer: otel.Tracer(tracing.Name()),
		gr:     converter.NewRepository[entities.UserDB, *entities.UserDB](db),
	}
}

func (repository *userRepository) Create(ctx context.Context, user *model.User) error {
	ctx, span := repository.tracer.Start(ctx, tracing.Name())
	defer span.End()

	var userDB entities.UserDB
	userDB.From(*user)

	result := repository.DB.WithContext(ctx).Create(&userDB)
	if err := result.Error; err != nil {
		return fmt.Errorf("Error creating user %v", err)
	}

	*user = userDB.Into()

	return nil
}

func (repository *userRepository) FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[model.User], error) {
	ctx, span := repository.tracer.Start(ctx, tracing.Name())
	defer span.End()

	return repository.gr.FindAllPaginated(ctx, pageable, "Permission.Roles")
}

func (repository *userRepository) FindByCode(ctx context.Context, code uuid.UUID) (*model.User, error) {
	ctx, span := repository.tracer.Start(ctx, tracing.Name())
	defer span.End()

	var userDB entities.UserDB
	result := repository.WithContext(ctx).Find(&userDB, "code = ?", code)

	if err := result.Error; err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("User not found")
	}

	user := userDB.Into()

	return &user, nil
}

func (repository *userRepository) FindByUsername(ctx context.Context, username string) (nilo.Option[model.User], error) {
	ctx, span := repository.tracer.Start(ctx, tracing.Name())
	defer span.End()

	var userDB entities.UserDB
	result := repository.WithContext(ctx).
		Preload("Permission.Roles").
		Find(&userDB, "username = ?", username)

	if err := result.Error; err != nil {
		return nilo.None[model.User](), err
	}

	if result.RowsAffected == 0 {
		return nilo.None[model.User](), nil
	}

	user := userDB.Into()

	return nilo.Some(user), nil
}
