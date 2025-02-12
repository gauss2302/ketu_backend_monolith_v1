package interfaces

import (
	"context"
	"ketu_backend_monolith_v1/internal/domain"
)

type OwnerRepository interface {
	Create(ctx context.Context, owner *domain.Owner) error
	GetByID(ctx context.Context, id uint) (*domain.Owner, error)
	GetByEmail(ctx context.Context, email string) (*domain.Owner, error)
	GetAll(ctx context.Context) ([]domain.Owner, error)
	Update(ctx context.Context, owner *domain.Owner) error
	Delete(ctx context.Context, id uint) error
}
