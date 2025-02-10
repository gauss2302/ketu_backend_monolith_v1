package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		err := c.Next()
		
		log.Printf(
			"Method: %s, Path: %s, Status: %d, Duration: %v, Error: %v",
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			time.Since(start),
			err,
		)
		
		return err
	}
} 