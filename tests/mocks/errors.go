package mocks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

// Mock backend.Error
type MockBackendError struct {
	mock.Mock
}

func (m *MockBackendError) ToResponse(c *fiber.Ctx) error {
	return nil
}
