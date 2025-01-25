// internal/handler/http/auth.go
package http

import (
	"errors"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/service"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequestDTO
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Body parsing failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload: " + err.Error(),
		})
	}
	log.Printf("Parsed registration request: %+v", req) // Log the parsed request

	response, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		log.Printf("Registration failed with error: %v", err) // Log the detailed error
		if errors.Is(err, domain.ErrEmailExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user: " + err.Error(), // Include error details
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Now h.validator will be defined
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	response, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		log.Printf("Login error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Refresh token is required",
        })
    }

    // Check Bearer token format
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid authorization header format",
        })
    }

    refreshToken := parts[1]
    response, err := h.authService.RefreshToken(c.Context(), refreshToken)
    if err != nil {
        if errors.Is(err, domain.ErrInvalidCredentials) {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid or expired refresh token",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to refresh token",
        })
    }

    return c.Status(fiber.StatusOK).JSON(response)
}
