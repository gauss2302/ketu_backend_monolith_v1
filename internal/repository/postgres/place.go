package postgres

import (
	"context"
	"database/sql"
	"ketu_backend_monolith_v1/internal/domain"
)

type PlaceRepo struct {
	db *sql.DB
}

func NewPlaceRepository(db *sql.DB) *PlaceRepo {
	return &PlaceRepo{
		db: db,
	}
}

func (r *PlaceRepo) CreatePlace(ctx context.Context, place *domain.Place) error {
	query := `
			INSERT INTO places (name, description, location, main_image, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		place.Name,
		place.Description,
		place.Location,
		place.MainImage,
		place.CreatedAt,
		place.UpdatedAt).Scan(&place.ID)
}
