package middleware

import (
	"fmt"
	configs "ketu_backend_monolith_v1/internal/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	config configs.JWTConfig
}

func NewAuthMiddleware(config configs.JWTConfig) *AuthMiddleware {
	return &AuthMiddleware{
		config: config,
	}
}

func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is missing",
				"code":  "AUTH_HEADER_MISSING",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format. Use 'Bearer <token>'",
				"code":  "INVALID_AUTH_FORMAT",
			})
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.config.AccessSecret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid token: %v", err),
				"code":  "INVALID_TOKEN",
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Store all relevant claims in context
			c.Locals("user_id", uint(claims["user_id"].(float64)))
			c.Locals("email", claims["email"].(string))
			if role, ok := claims["role"].(string); ok {
				c.Locals("role", role)
			}
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}
}
