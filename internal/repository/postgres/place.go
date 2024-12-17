package postgres

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *PlaceRepo) GetPlaceByID(ctx context.Context, id uint) (*domain.Place, error) {
	query := `
			SELECT id, name, description, location, main_image, created_at, updated_at
			FROM places
			WHERE id = $1`

	var place domain.Place
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&place.ID, &place.Name, &place.Description, &place.Location, &place.MainImage, &place.CreatedAt, &place.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &place, nil
}

func (r *PlaceRepo) UpdatePlaceByID(ctx context.Context, place *domain.Place) error {
	query := `
	UPDATE places SET name = $1, description = $2, location = $3, main_image = $4, updated_at = $5
		WHERE id = $6`

	result, err := r.db.ExecContext(ctx, query,
		place.Name,
		place.Description,
		place.Location,
		place.MainImage,
		place.UpdatedAt,
		place.ID)
	if err != nil {
		return fmt.Errorf("error updating place: %w", err)
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrPlaceNotFound
	}
	return nil
}

func (r *PlaceRepo) DeletePlaceByID(ctx context.Context, id uint) error {
	query := `DELETE FROM places WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("error deleting place: %w", err)
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrPlaceNotFound
	}

	return nil
}
