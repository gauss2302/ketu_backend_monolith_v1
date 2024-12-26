package http

import (
	"errors"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"

	"ketu_backend_monolith_v1/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PlaceHandler struct {
	placeService *service.PlaceService
}

func NewPlaceHandler(placeService *service.PlaceService) *PlaceHandler {
	return &PlaceHandler{
		placeService: placeService,
	}
}

// Create handles place creation
func (h *PlaceHandler) Create(c *fiber.Ctx) error {
	var createDTO dto.PlaceCreateDTO
	if err := c.BodyParser(&createDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	response, err := h.placeService.CreatePlace(c.Context(), &createDTO)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		case errors.Is(err, domain.ErrPlaceExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Place already exists",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create place",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetByID handles retrieving a place by ID
func (h *PlaceHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid place ID",
		})
	}

	response, err := h.placeService.GetPlace(c.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPlaceNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Place not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve place",
			})
		}
	}

	return c.JSON(response)
}

// List handles retrieving a paginated list of places
func (h *PlaceHandler) List(c *fiber.Ctx) error {
	params := dto.PaginationParams{
		Offset: getIntQuery(c, "offset", 0),
		Limit:  getIntQuery(c, "limit", 10),
	}

	response, err := h.placeService.ListPlaces(c.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidPagination):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid pagination parameters",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve places",
			})
		}
	}

	return c.JSON(response)
}

// Search handles searching places with criteria
func (h *PlaceHandler) Search(c *fiber.Ctx) error {
	searchDTO := dto.PlaceSearchDTO{
		PaginationParams: dto.PaginationParams{
			Offset: getIntQuery(c, "offset", 0),
			Limit:  getIntQuery(c, "limit", 10),
		},
		Name:      c.Query("name"),
		City:      c.Query("city"),
		Province:  c.Query("province"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order", "asc"),
	}

	response, err := h.placeService.SearchPlaces(c.Context(), &searchDTO)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCriteria):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid search criteria",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to search places",
			})
		}
	}

	return c.JSON(response)
}

// Update handles updating an existing place
func (h *PlaceHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid place ID",
		})
	}

	var updateDTO dto.PlaceUpdateDTO
	if err := c.BodyParser(&updateDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	response, err := h.placeService.UpdatePlace(c.Context(), uint(id), &updateDTO)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPlaceNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Place not found",
			})
		case errors.Is(err, domain.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update place",
			})
		}
	}

	return c.JSON(response)
}

// Delete handles removing a place
func (h *PlaceHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid place ID",
		})
	}

	err = h.placeService.DeletePlace(c.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPlaceNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Place not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete place",
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
