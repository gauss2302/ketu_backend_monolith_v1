package http

import (
	"github.com/gofiber/fiber/v2"
)

func CreateHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Handle(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "healthy",
	})
}