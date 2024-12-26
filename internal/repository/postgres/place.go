// postgres/place.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	repository "ketu_backend_monolith_v1/internal/repository/interfaces"

	"github.com/jmoiron/sqlx"
)

type placeRepository struct {
	db *sqlx.DB
}

// NewPlaceRepository creates a new instance of PlaceRepository
func NewPlaceRepository(db *sqlx.DB) repository.PlaceRepository {
	if db == nil {
		panic("nil db provided to NewPlaceRepository")
	}
	return &placeRepository{
		db: db,
	}
}

func (r *placeRepository) Create(ctx context.Context, place *domain.Place) error {
	locationJSON, err := json.Marshal(place.Location)
	if err != nil {
		return fmt.Errorf("error marshaling location: %w", err)
	}

	query := `
        INSERT INTO places (name, description, location, main_image, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

	err = r.db.QueryRowContext(ctx, query,
		place.Name,
		place.Description,
		locationJSON,
		place.MainImage,
		place.CreatedAt,
		place.UpdatedAt).Scan(&place.ID)

	if err != nil {
		return fmt.Errorf("error creating place: %w", err)
	}

	return nil
}

func (r *placeRepository) List(ctx context.Context, offset, limit int) ([]*domain.Place, int, error) {
	// First, get total count
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM places").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting places: %w", err)
	}

	// Then get paginated results
	query := `
        SELECT p.id, p.name, p.description, p.location, p.main_image, 
               p.created_at, p.updated_at,
               COALESCE(array_agg(pi.image_url) FILTER (WHERE pi.image_url IS NOT NULL), '{}') as images
        FROM places p
        LEFT JOIN place_images pi ON p.id = pi.place_id
        GROUP BY p.id
        ORDER BY p.created_at DESC
        LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing places: %w", err)
	}
	defer rows.Close()

	var places []*domain.Place
	for rows.Next() {
		var place domain.Place
		var locationJSON json.RawMessage

		err := rows.Scan(
			&place.ID,
			&place.Name,
			&place.Description,
			&locationJSON,
			&place.MainImage,
			&place.CreatedAt,
			&place.UpdatedAt,
			&place.Images,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning place: %w", err)
		}

		if err := json.Unmarshal(locationJSON, &place.Location); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling location: %w", err)
		}

		places = append(places, &place)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating places: %w", err)
	}

	return places, total, nil
}

func (r *placeRepository) Search(ctx context.Context, criteria dto.PlaceSearchCriteria) ([]*domain.Place, int, error) {
	// Validate criteria
	if criteria.Limit <= 0 {
		criteria.Limit = 10 // default limit
	}
	if criteria.Offset < 0 {
		criteria.Offset = 0
	}

	// Build the WHERE clause dynamically
	where := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if criteria.Name != "" {
		where += fmt.Sprintf(" AND LOWER(name) LIKE LOWER($%d)", argCount)
		args = append(args, "%"+criteria.Name+"%")
		argCount++
	}

	if criteria.City != "" {
		where += fmt.Sprintf(" AND location->>'city' = $%d", argCount)
		args = append(args, criteria.City)
		argCount++
	}

	if criteria.Province != "" {
		where += fmt.Sprintf(" AND location->>'province' = $%d", argCount)
		args = append(args, criteria.Province)
		argCount++
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM places " + where
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting search results: %w", err)
	}

	// Build the final query
	orderBy := " ORDER BY created_at DESC"
	if criteria.SortBy != "" {
		orderBy = fmt.Sprintf(" ORDER BY %s %s", criteria.SortBy, criteria.SortOrder)
	}

	query := fmt.Sprintf(`
        SELECT p.id, p.name, p.description, p.location, p.main_image, 
               p.created_at, p.updated_at,
               COALESCE(array_agg(pi.image_url) FILTER (WHERE pi.image_url IS NOT NULL), '{}') as images
        FROM places p
        LEFT JOIN place_images pi ON p.id = pi.place_id
        %s
        GROUP BY p.id
        %s
        LIMIT $%d OFFSET $%d`,
		where, orderBy, argCount, argCount+1)

	args = append(args, criteria.Limit, criteria.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error searching places: %w", err)
	}
	defer rows.Close()

	var places []*domain.Place
	for rows.Next() {
		var place domain.Place
		var locationJSON json.RawMessage

		err := rows.Scan(
			&place.ID,
			&place.Name,
			&place.Description,
			&locationJSON,
			&place.MainImage,
			&place.CreatedAt,
			&place.UpdatedAt,
			&place.Images,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning place: %w", err)
		}

		if err := json.Unmarshal(locationJSON, &place.Location); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling location: %w", err)
		}

		places = append(places, &place)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating search results: %w", err)
	}

	return places, total, nil
}

func (r *placeRepository) GetByID(ctx context.Context, id uint) (*domain.Place, error) {
	query := `
		 SELECT p.id, p.name, p.description, p.location, p.main_image, 
				  p.created_at, p.updated_at,
				  COALESCE(array_agg(pi.image_url) FILTER (WHERE pi.image_url IS NOT NULL), '{}') as images
		 FROM places p
		 LEFT JOIN place_images pi ON p.id = pi.place_id
		 WHERE p.id = $1
		 GROUP BY p.id`

	var place domain.Place
	var locationJSON json.RawMessage

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&place.ID,
		&place.Name,
		&place.Description,
		&locationJSON,
		&place.MainImage,
		&place.CreatedAt,
		&place.UpdatedAt,
		&place.Images,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPlaceNotFound
		}
		return nil, fmt.Errorf("error getting place by ID: %w", err)
	}

	if err := json.Unmarshal(locationJSON, &place.Location); err != nil {
		return nil, fmt.Errorf("error unmarshaling location: %w", err)
	}

	return &place, nil
}

func (r *placeRepository) Update(ctx context.Context, place *domain.Place) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	locationJSON, err := json.Marshal(place.Location)
	if err != nil {
		return fmt.Errorf("error marshaling location: %w", err)
	}

	// Update main place record
	query := `
		 UPDATE places 
		 SET name = $1, 
			  description = $2, 
			  location = $3, 
			  main_image = $4, 
			  updated_at = $5
		 WHERE id = $6`

	result, err := tx.ExecContext(ctx, query,
		place.Name,
		place.Description,
		locationJSON,
		place.MainImage,
		place.UpdatedAt,
		place.ID,
	)
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

	// Delete existing images
	_, err = tx.ExecContext(ctx, "DELETE FROM place_images WHERE place_id = $1", place.ID)
	if err != nil {
		return fmt.Errorf("error deleting existing images: %w", err)
	}

	// Insert new images if present
	if len(place.Images) > 0 {
		imageQuery := `
			  INSERT INTO place_images (place_id, image_url)
			  VALUES ($1, $2)`

		for _, img := range place.Images {
			_, err := tx.ExecContext(ctx, imageQuery, place.ID, img)
			if err != nil {
				return fmt.Errorf("error inserting place image: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (r *placeRepository) Delete(ctx context.Context, id uint) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete associated images first
	_, err = tx.ExecContext(ctx, "DELETE FROM place_images WHERE place_id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting place images: %w", err)
	}

	// Delete the place
	result, err := tx.ExecContext(ctx, "DELETE FROM places WHERE id = $1", id)
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

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
