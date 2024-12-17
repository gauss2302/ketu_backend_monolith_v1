package service

import (
	"context"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/repository"
)

type PlaceUseCase struct {
	placeRepo repository.PlaceRepository
}

func NewPlaceUsecase(placeRepo repository.PlaceRepository) *PlaceUseCase {
	return &PlaceUseCase{placeRepo: placeRepo}
}

func (p *PlaceUseCase) CreatePlace(ctx context.Context, place *domain.Place) error {
	return p.placeRepo.Create(ctx, place)
}

func (p *PlaceUseCase) GetPlaceByID(ctx context.Context, id uint) (*domain.Place, error) {
	return p.placeRepo.GetByID(ctx, id)
}

func (p *PlaceUseCase) UpdatePlaceByID(ctx context.Context, id uint, place *domain.Place) error {
	return p.placeRepo.UpdateByID(ctx, id, place)
}

func (p *PlaceUseCase) DeletePlaceByID(ctx context.Context, id uint) error {
	return p.placeRepo.DeletePlaceByID(ctx, id)
}

func (p *PlaceUseCase) GetAllPlaces(ctx context.Context) ([]domain.Place, error) {
	return p.placeRepo.GetAll(ctx)
}

func (p *PlaceUseCase) GetPlaceByName(ctx context.Context, name string) (*domain.Place, error) {
	return p.placeRepo.GetbyName(ctx, name)
}

func (p *PlaceUseCase) GetPlaceByLocation(ctx context.Context, location string) (*domain.Place, error) {
	return p.placeRepo.GetbyLocation(ctx, location)
}
