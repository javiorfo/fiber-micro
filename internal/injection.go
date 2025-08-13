package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/javiorfo/fiber-micro/adapter/database"
	"github.com/javiorfo/fiber-micro/adapter/database/repository"
	"github.com/javiorfo/fiber-micro/adapter/http/routes"
	"github.com/javiorfo/fiber-micro/application/domain/service"
	"github.com/javiorfo/go-microservice-lib/security"
)

func InjectDependencies(api fiber.Router) {
	db := database.DBinstance
	tokenSecurity := security.NewTokenSecurity()

	userRepository := repository.NewUserRepository(db)
	permRepository := repository.NewPermissionRepository(db)
	userService := service.NewUserService(userRepository, permRepository)

	routes.User(api, tokenSecurity, userService)
}
