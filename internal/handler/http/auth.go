// internal/handler/http/auth.go
package http

import (
	"errors"
	"fmt"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/service"
	"log"

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
	log.Printf("Handling register request")
	
	req, ok := c.Locals("validated").(*dto.RegisterRequestDTO)
	if !ok {
		log.Printf("Failed to get validated request from context")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("Processing registration for email: %s", req.Email)

	response, err := h.authService.Register(c.Context(), req)
	if err != nil {
		log.Printf("Registration error: %v", err)
		if errors.Is(err, domain.ErrEmailExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to register user: %v", err),
		})
	}

	log.Printf("Registration successful for email: %s", req.Email)
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
    userID := c.Locals("user_id").(uint)
    
    response, err := h.authService.RefreshToken(c.Context(), userID)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Failed to refresh token",
            "details": err.Error(),
        })
    }

    return c.JSON(response)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    
    if err := h.authService.Logout(c.Context(), userID); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to logout",
        })
    }

    return c.SendStatus(fiber.StatusNoContent)
}
