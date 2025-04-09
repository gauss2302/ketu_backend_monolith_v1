package http

import (
	"errors"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RestaurantHandler struct {
	restaurantService *service.RestaurantService
}

func NewRestaurantHandler(restaurantService *service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
	}
}

// Create handles restaurant creation
func (h *RestaurantHandler) Create(c *fiber.Ctx) error {
	var createDTO dto.CreateRestaurantDTO
	if err := c.BodyParser(&createDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Get owner ID from authenticated user
	ownerID := c.Locals("user_id").(uint)
	createDTO.OwnerID = ownerID

	response, err := h.restaurantService.Create(c.Context(), &createDTO)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create restaurant",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetByID handles retrieving a restaurant by ID
func (h *RestaurantHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid restaurant ID",
		})
	}

	response, err := h.restaurantService.GetByID(c.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRestaurantNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Restaurant not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve restaurant",
			})
		}
	}

	return c.JSON(response)
}

// List handles retrieving a paginated list of restaurants
func (h *RestaurantHandler) List(c *fiber.Ctx) error {
	params := dto.ListParams{
		Offset: getIntQuery(c, "offset", 0),
		Limit:  getIntQuery(c, "limit", 10),
	}

	restaurants, total, err := h.restaurantService.List(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve restaurants",
		})
	}

	return c.JSON(fiber.Map{
		"data": restaurants,
		"pagination": dto.ListResponse{
			Total:  total,
			Offset: params.Offset,
			Limit:  params.Limit,
		},
	})
}

// ListByOwner handles retrieving restaurants for the authenticated owner
func (h *RestaurantHandler) ListByOwner(c *fiber.Ctx) error {
	ownerID := c.Locals("owner_id").(uint)
	params := dto.ListParams{
		Offset: getIntQuery(c, "offset", 0),
		Limit:  getIntQuery(c, "limit", 10),
	}

	restaurants, total, err := h.restaurantService.ListByOwnerID(c.Context(), ownerID, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve owner's restaurants",
		})
	}

	return c.JSON(fiber.Map{
		"data": restaurants,
		"pagination": dto.ListResponse{
			Total:  total,
			Offset: params.Offset,
			Limit:  params.Limit,
		},
	})
}

// Update handles updating an existing restaurant
func (h *RestaurantHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid restaurant ID",
		})
	}

	var updateDTO dto.UpdateRestaurantDTO
	if err := c.BodyParser(&updateDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Verify ownership
	ownerID := c.Locals("user_id").(uint)

	response, err := h.restaurantService.Update(c.Context(), uint(id), ownerID, &updateDTO)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRestaurantNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Restaurant not found",
			})
		case errors.Is(err, domain.ErrUnauthorized):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Not authorized to update this restaurant",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update restaurant",
			})
		}
	}

	return c.JSON(response)
}

// Delete handles removing a restaurant
func (h *RestaurantHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid restaurant ID",
		})
	}

	// Verify ownership
	ownerID := c.Locals("owner_id").(uint)

	err = h.restaurantService.Delete(c.Context(), uint(id), ownerID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRestaurantNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Restaurant not found",
			})
		case errors.Is(err, domain.ErrUnauthorized):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Not authorized to delete this restaurant",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete restaurant",
			})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Helper function to get query parameters with default values
func getIntQuery(c *fiber.Ctx, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
