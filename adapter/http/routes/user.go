package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/javiorfo/fiber-micro/adapter/http/handlers"
	"github.com/javiorfo/fiber-micro/application/port"
	"github.com/javiorfo/go-microservice-lib/security"
)

func User(app fiber.Router, auth security.Authorizer, service port.UserService) {
	app.Get("/users", auth.Secure("ROLE_1", "ROLE_2"), handlers.FindAllUsers(service))
	app.Post("/users", auth.Secure("ROLE_1"), handlers.CreateUser(service))
	app.Post("/users/login", handlers.Login(service))
}
