package http

import (
	"github.com/gofiber/fiber/v2"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/service"
)

type PlaceHandler struct {
	placeUsecase *service.PlaceUseCase
}

func NewPlaceHandler(placeUsecase *service.PlaceUseCase) *PlaceHandler {
	return &PlaceHandler{
		placeUsecase: placeUsecase,
	}
}

func (h *PlaceHandler) Create(c *fiber.Ctx) error {
	var place domain.Place
	if err := c.BodyParser(&place); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if err := h.placeUsecase.CreatePlace(c.Context(), &place); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot create place",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(place)
}
