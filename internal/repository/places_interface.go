// user interface

package repository

import (
	"context"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/handler/dto"
)

// PlaceRepository defines the interface for place data operations
type PlaceRepository interface {
	// Create stores a new place in the repository
	Create(ctx context.Context, place *domain.Place) error

	// GetByID retrieves a place by its ID
	GetByID(ctx context.Context, id uint) (*domain.Place, error)

	// Update updates an existing place
	Update(ctx context.Context, place *domain.Place) error

	// Delete removes a place by its ID
	Delete(ctx context.Context, id uint) error

	// List retrieves a paginated list of places
	List(ctx context.Context, offset, limit int) ([]*domain.Place, int, error)

	// Search retrieves places based on search criteria
	Search(ctx context.Context, criteria dto.PlaceSearchCriteria) ([]*domain.Place, int, error)
}
