package repository

import (
	"context"
	"fmt"

	"github.com/javiorfo/fiber-micro/adapter/database/entities"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/fiber-micro/application/port"
	"github.com/javiorfo/go-microservice-lib/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type roleRepository struct {
	*gorm.DB
	tracer trace.Tracer
}

func NewRoleRepository(db *gorm.DB) port.RoleRepository {
	return &roleRepository{DB: db, tracer: otel.Tracer(tracing.Name())}
}

func (repository *roleRepository) Create(ctx context.Context, role *model.Role) error {
	ctx, span := repository.tracer.Start(ctx, tracing.Name())
	defer span.End()

	var roleDB entities.RoleDB
	roleDB.From(*role)

	result := repository.DB.WithContext(ctx).Create(&roleDB)
	if err := result.Error; err != nil {
		return fmt.Errorf("Error creating role %v", err)
	}

	*role = roleDB.Into()

	return nil
}
