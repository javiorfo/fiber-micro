package injection

import (
	"github.com/gofiber/fiber/v2"
	"github.com/javiorfo/fiber-micro/adapter/database/connection"
	"github.com/javiorfo/fiber-micro/adapter/database/repository"
)

func InjectDependencies(api fiber.Router) {
	db := connection.DBinstance

	userRepository := repository.NewUserRepository(db)
	_ = userRepository
}
