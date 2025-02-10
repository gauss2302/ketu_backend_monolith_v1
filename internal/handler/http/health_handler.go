package http

import "github.com/gofiber/fiber/v2"

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Renamed method to avoid conflict
func (h *HealthHandler) HandleHealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
