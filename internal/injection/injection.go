package injection

import (
	"github.com/gofiber/fiber/v2"
	"github.com/javiorfo/fiber-micro/adapter/database/connection"
)

func InjectDependencies(api fiber.Router) {
	db := connection.DBinstance
	_ = db
}
