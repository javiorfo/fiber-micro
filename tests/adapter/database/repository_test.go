package database

import (
	"context"
	"testing"

	"github.com/javiorfo/fiber-micro/adapter/database/entities"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/go-microservice-lib/pagination"
)

func TestDatabase(t *testing.T) {
	ctx := context.Background()

	role1 := model.NewRole("ROLE1")
	if err := roleRepo.Create(ctx, role1); err != nil {
		t.Fatalf("Failed to insert role1: %v", err)
	}

	role2 := model.NewRole("ROLE2")
	if err := roleRepo.Create(ctx, role2); err != nil {
		t.Fatalf("Failed to insert role2: %v", err)
	}

	permission := model.NewPermission("PERM", []model.Role{*role1, *role2})
	if err := permRepo.Create(ctx, permission); err != nil {
		t.Fatalf("Failed to insert permission: %v", err)
	}

	password := "1234"
	user := model.NewUser("Javi", "javi@mail.com", *permission, password, "auditor")

	if err := userRepo.Create(ctx, &user); err != nil {
		t.Fatalf("Failed to insert record: %v", err)
	}

	userResult, err := userRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		t.Fatalf("Failed to query record: %v", err)
	}

	if userResult.CreatedBy != "auditor" {
		t.Errorf("Expected name to be 'auditor', got '%s'", userResult.CreatedBy)
	}

	if userResult.ID != 1 {
		t.Errorf("Expected name to be '1', got '%d'", userResult.ID)
	}

	// User inserts to test pagination
	user2 := model.NewUser("Carlos", "carlos@mail.com", *permission, password, "auditor")
	if err := userRepo.Create(ctx, &user2); err != nil {
		t.Fatalf("Failed to insert record: %v", err)
	}

	user3 := model.NewUser("Lionel", "lionel@mail.com", *permission, password, "auditor")
	if err := userRepo.Create(ctx, &user3); err != nil {
		t.Fatalf("Failed to insert record: %v", err)
	}

	users, err := userRepo.FindAll(ctx, entities.NewUserFilter(pagination.DefaultPage(), "Javi", "PERM", ""))
	if err != nil {
		t.Fatalf("Failed to execute findAll: %v", err)
	}

	if len(users) == 0 {
		t.Error("Expected '1', got '0'")
	}
}
