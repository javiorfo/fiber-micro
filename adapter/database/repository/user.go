package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/javiorfo/fiber-micro/adapter/database/entities"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/fiber-micro/application/port"
	"github.com/javiorfo/go-microservice-lib/pagination"
	"github.com/javiorfo/go-microservice-lib/tracing"
	"github.com/javiorfo/steams"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type userRepository struct {
	*gorm.DB
	tracer trace.Tracer
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &userRepository{DB: db, tracer: otel.Tracer(tracing.Name())}
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

func (repository *userRepository) FindAll(ctx context.Context, queryFilter pagination.QueryFilter) ([]model.User, error) {
	ctx, span := repository.tracer.Start(ctx, tracing.Name())
	defer span.End()

	var usersDB []entities.UserDB
	filter := repository.WithContext(ctx).
		Preload("Permission.Roles").
		Joins("INNER JOIN permissions ON users.permission_id = permissions.id")

	filter, err := queryFilter.Filter(filter)
	if err != nil {
		return nil, err
	}

	results := filter.Find(&usersDB)

	if err := results.Error; err != nil {
		return nil, err
	}

	users := steams.Mapping(steams.OfSlice(usersDB), func(userDB entities.UserDB) model.User {
		return userDB.Into()
	}).Collect()

	return users, nil
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

func (repository *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	ctx, span := repository.tracer.Start(ctx, tracing.Name())
	defer span.End()

	var userDB entities.UserDB
	result := repository.WithContext(ctx).
		Preload("Permission.Roles").
		Find(&userDB, "username = ?", username)

	if err := result.Error; err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("User not found")
	}

	user := userDB.Into()

	return &user, nil
}
