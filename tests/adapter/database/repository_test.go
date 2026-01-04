package database

import (
	"context"
	"testing"

	"github.com/javiorfo/fiber-micro/adapter/database/entities"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/gormen/pagination"
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

	permission := model.NewPermission("PERM2", []model.Role{*role1, *role2})
	if err := permRepo.Create(ctx, permission); err != nil {
		t.Fatalf("Failed to insert permission: %v", err)
	}

	password := "1234"
	user := model.NewUser("Javi", "javi@mail.com", *permission, password, "auditor")

	if err := userRepo.Create(ctx, &user); err != nil {
		t.Fatalf("Failed to insert record: %v", err)
	}

	// Find by Username
	userResult, err := userRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		t.Fatalf("Failed to query record: %v", err)
	}

	if userResult.Unwrap().CreatedBy != "auditor" {
		t.Errorf("Expected name to be 'auditor', got '%s'", userResult.Unwrap().CreatedBy)
	}

	if userResult.Unwrap().ID != 3 {
		t.Errorf("Expected name to be '3', got '%d'", userResult.Unwrap().ID)
	}

	if userResult.Unwrap().Permission.Name != "PERM2" {
		t.Errorf("Expected Permission Name to be 'PERM2', got '%v'", userResult.Unwrap().Permission)
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

	// Find all
	pageRequest, err := pagination.PageRequestFrom(1, 10, pagination.WithFilter(entities.NewUserFilter("Javi", "PERM2", "")))
	page, err := userRepo.FindAll(ctx, pageRequest)
	if err != nil {
		t.Fatalf("Failed to execute findAll: %v", err)
	}

	if len(page.Elements) == 0 {
		t.Error("Expected '1', got '0'")
	}

	if page.Total == 0 {
		t.Error("Expected '1', got '0'")
	}

	// Diff filter and page against count
	pageRequest, err = pagination.PageRequestFrom(1, 1, pagination.WithFilter(entities.NewUserFilter("", "", "")))
	page, err = userRepo.FindAll(ctx, pageRequest)
	if err != nil {
		t.Fatalf("Failed to execute findAll: %v", err)
	}

	if len(page.Elements) != 1 {
		t.Errorf("Expected '1', got %d", len(page.Elements))
	}

	if page.Total != 5 {
		t.Errorf("Expected '5', got %d", page.Total)
	}
}
