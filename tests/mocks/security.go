package mocks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

// Mock Security
type MockAuthorizer struct {
	mock.Mock
}

func (m *MockAuthorizer) Secure(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
