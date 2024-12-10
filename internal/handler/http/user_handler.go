package http

import (
	"github.com/gofiber/fiber/v2"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/handler/dto"
	"ketu_backend_monolith_v1/internal/service"
)

type UserHandler struct {
	userUsecase *service.UserUseCase
}

func NewUserHandler(userService *service.UserUseCase) *UserHandler {
	return &UserHandler{
		userUsecase: userService,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var input dto.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	user, err := h.userUsecase.CreateUser(c.Context(), input)
	if err != nil {
		if err == domain.ErrEmailExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}
