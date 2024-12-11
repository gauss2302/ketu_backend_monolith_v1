package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
)

var validate = validator.New()

func ValidateBody(payload interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("Validating request body for path: %s", c.Path())

		if err := c.BodyParser(payload); err != nil {
			log.Printf("Body parsing error: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request payload",
			})
		}

		if err := validate.Struct(payload); err != nil {
			errors := []map[string]string{}
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, map[string]string{
					"field": err.Field(),
					"error": err.Tag(),
				})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": errors,
			})
		}

		c.Locals("validated", payload)
		return c.Next()
	}
}
