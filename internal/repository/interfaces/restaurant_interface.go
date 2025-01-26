package interfaces

import (
	"context"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
)

type RestaurantRepository interface {
	Create(ctx context.Context, restaurant *domain.Restaurant) error
	GetByID(ctx context.Context, id uint) (*domain.Restaurant, error)
	List(ctx context.Context, params dto.ListParams) ([]domain.Restaurant, int, error)
	Update(ctx context.Context, restaurant *domain.Restaurant) error
	Delete(ctx context.Context, id uint) error
	ListByOwnerID(ctx context.Context, ownerID uint, params dto.ListParams) ([]domain.Restaurant, int, error)
} 