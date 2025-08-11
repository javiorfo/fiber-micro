package database

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/javiorfo/fiber-micro/adapter/database/migrator/tables"
	"github.com/javiorfo/fiber-micro/adapter/database/repository"
	"github.com/javiorfo/fiber-micro/application/port"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var container testcontainers.Container
var userRepo port.UserRepository
var roleRepo port.RoleRepository
var permRepo port.PermissionRepository

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	var err error
	container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start container: %s", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %s", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Failed to get container port: %s", err)
	}

	dsn := "host=" + host + " port=" + port.Port() + " user=testuser password=testpass dbname=testdb sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	tables.Migrate(db, "../../../migrations")

	userRepo = repository.NewUserRepository(db)
	roleRepo = repository.NewRoleRepository(db)
	permRepo = repository.NewPermissionRepository(db)

	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("Failed to terminate container: %s", err)
	}

	os.Exit(code)
}
