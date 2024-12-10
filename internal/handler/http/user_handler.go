package http

import (
	"github.com/gofiber/fiber/v2"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/handler/dto"
	"ketu_backend_monolith_v1/internal/service"
	"strconv"
)

type UserHandler struct {
	userService *service.UserUseCase
}

func NewUserHandler(userService *service.UserUseCase) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var input dto.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	user, err := h.userService.CreateUser(c.Context(), input)
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

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.userService.GetUser(c.Context(), uint(id))
	if err != nil {
		if err == domain.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	return c.JSON(dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	users, err := h.userService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get users",
		})
	}

	var response []dto.UserResponse
	for _, user := range users {
		response = append(response, dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return c.JSON(response)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var input dto.UpdateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := h.userService.UpdateUser(c.Context(), uint(id), input); err != nil {
		if err == domain.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := h.userService.DeleteUser(c.Context(), uint(id)); err != nil {
		if err == domain.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
