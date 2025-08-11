package injection

import (
	"github.com/gofiber/fiber/v2"
	"github.com/javiorfo/fiber-micro/adapter/database/connection"
	"github.com/javiorfo/fiber-micro/adapter/database/repository"
	"github.com/javiorfo/fiber-micro/application/domain/service"
)

func InjectDependencies(api fiber.Router) {
	db := connection.DBinstance

	userRepository := repository.NewUserRepository(db)
	permRepository := repository.NewPermissionRepository(db)
	userService := service.NewUserService(userRepository, permRepository)
	_ = userService
}
