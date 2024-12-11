package http

import (
	"github.com/gofiber/fiber/v2"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/handler/dto"
	"ketu_backend_monolith_v1/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input dto.RegisterRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	user, err := h.authService.Register(c.Context(), input)
	if err != nil {
		if err == domain.ErrEmailExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input dto.LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	user, accessToken, refreshToken, err := h.authService.Login(c.Context(), input)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}

	return c.JSON(dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	})
}
