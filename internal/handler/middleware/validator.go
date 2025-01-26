package middleware

import (
	"log"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateBody(payload interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("Validating request body for path: %s", c.Path())

		// Create a new instance of the payload
		p := reflect.New(reflect.TypeOf(payload).Elem()).Interface()

		if err := c.BodyParser(p); err != nil {
			log.Printf("Body parsing error: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
				"details": err.Error(),
			})
		}

		if err := validate.Struct(p); err != nil {
			log.Printf("Validation error: %v", err)
			var errors []map[string]string
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

		c.Locals("validated", p)
		return c.Next()
	}
}
