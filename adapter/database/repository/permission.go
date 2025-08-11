package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/javiorfo/fiber-micro/adapter/database/entities"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/fiber-micro/application/port"
	"github.com/javiorfo/go-microservice-lib/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type permissionRepository struct {
	*gorm.DB
	tracer trace.Tracer
}

func NewPermissionRepository(db *gorm.DB) port.PermissionRepository {
	return &permissionRepository{DB: db, tracer: otel.Tracer(tracing.Name())}
}

func (repository *permissionRepository) Create(ctx context.Context, perm *model.Permission) error {
	ctx, span := repository.tracer.Start(ctx, tracing.Name())
	defer span.End()

	var permDB entities.PermissionDB
	permDB.From(*perm)

	result := repository.DB.WithContext(ctx).Create(&permDB)
	if err := result.Error; err != nil {
		return fmt.Errorf("Error creating permission %v", err)
	}

	*perm = permDB.Into()

	return nil
}

func (repository *permissionRepository) FindByName(ctx context.Context, name string) (*model.Permission, error) {
	ctx, span := repository.tracer.Start(ctx, tracing.Name())
	defer span.End()

	var permDB entities.PermissionDB
	result := repository.WithContext(ctx).Find(&permDB, "name = ?", name)

	if err := result.Error; err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("User not found")
	}

	permission := permDB.Into()

	return &permission, nil
}
