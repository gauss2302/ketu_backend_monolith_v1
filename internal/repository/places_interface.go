package repository

import (
	"context"
	"ketu_backend_monolith_v1/internal/domain"
)

type PlaceRepository interface {
	Create(ctx context.Context, place *domain.Place) error
	GetByID(ctx context.Context, id uint) (*domain.Place, error)
	UpdateByID(ctx context.Context, id uint, place *domain.Place) error
	DeleteByID(ctx context.Context, id uint) error
	GetAll(ctx context.Context) ([]domain.Place, error)
	GetbyName(ctx context.Context, name string) (*domain.Place, error)
	GetbyLocation(ctx context.Context, location string) (*domain.Place, error)
}
