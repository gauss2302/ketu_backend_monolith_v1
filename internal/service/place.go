package service

import "ketu_backend_monolith_v1/internal/repository"

type PlaceUseCase struct {
	placeRepo repository.PlaceRepository
}

func NewPlaceUsecase(placeRepo repository.PlaceRepository) *PlaceUseCase {
	return &PlaceUseCase{placeRepo: placeRepo}
}
