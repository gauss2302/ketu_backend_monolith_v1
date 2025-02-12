package http

import (
	"database/sql"
	"errors"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/service"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OwnerAuthHandler struct {
	ownerAuthService *service.OwnerAuthService
	validator        *validator.Validate
}

func NewOwnerAuthHandler(ownerAuthService *service.OwnerAuthService) *OwnerAuthHandler {
	return &OwnerAuthHandler{
		ownerAuthService: ownerAuthService,
		validator:        validator.New(),
	}
}

func (h *OwnerAuthHandler) Register(c *fiber.Ctx) error {
	req, ok := c.Locals("validated").(*dto.OwnerRegisterRequestDTO)
	if !ok {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			  "error": "Invalid request body",
		 })
	}

	// Validate the request
	if err := h.validator.Struct(req); err != nil {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			  "error": "Validation failed",
			  "details": err.Error(),
		 })
	}

	response, err := h.ownerAuthService.Register(c.Context(), req)
	if err != nil {
		 if errors.Is(err, domain.ErrEmailExists) {
			  return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Email already exists",
			  })
		 }
		 return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			  "error": "Failed to register owner",
		 })
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *OwnerAuthHandler) Login(c *fiber.Ctx) error {
	req, ok := c.Locals("validated").(*dto.OwnerLoginRequestDTO)
	if !ok {
		log.Printf("Login validation error: request body is not of type OwnerLoginRequestDTO")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Log the request data (be careful with passwords in production)
	log.Printf("Login attempt for email: %s", req.Email)

	response, err := h.ownerAuthService.Login(c.Context(), req)
	if err != nil {
		log.Printf("Login error details: %+v", err)
		
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Owner not found",
				"details": "No owner account found with this email",
			})
		}
		
		// Log the full error for debugging
		log.Printf("Unexpected error during login: %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to login",
			"details": err.Error(),
			"message": "An unexpected error occurred during login",
		})
	}

	log.Printf("Login successful for email: %s", req.Email)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *OwnerAuthHandler) RefreshToken(c *fiber.Ctx) error {
	ownerID := c.Locals("user_id").(uint)

	response, err := h.ownerAuthService.RefreshToken(c.Context(), ownerID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Failed to refresh token",
			"details": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *OwnerAuthHandler) Logout(c *fiber.Ctx) error {
	ownerID := c.Locals("user_id").(uint)

	if err := h.ownerAuthService.Logout(c.Context(), ownerID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to logout",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
} 