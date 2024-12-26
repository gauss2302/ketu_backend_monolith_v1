package repository

import (
	"context"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/handler/dto"
)

// PlaceRepository defines the interface for place data operations
type PlaceRepository interface {
	Create(ctx context.Context, place *domain.Place) error
	GetByID(ctx context.Context, id uint) (*domain.Place, error)
	Update(ctx context.Context, place *domain.Place) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*domain.Place, int, error)
	Search(ctx context.Context, criteria dto.PlaceSearchCriteria) ([]*domain.Place, int, error)
}
