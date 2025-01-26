package postgres

import (
	"context"
	"fmt"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"

	"github.com/jmoiron/sqlx"
)

type RestaurantRepository struct {
	db *sqlx.DB
}

func NewRestaurantRepository(db *sqlx.DB) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

func (r *RestaurantRepository) Create(ctx context.Context, restaurant *domain.Restaurant) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert restaurant
	query := `
		INSERT INTO restaurants (owner_id, name, description, main_image, is_verified)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING restaurant_id`
	
	err = tx.QueryRowxContext(ctx, query,
		restaurant.OwnerID,
		restaurant.Name,
		restaurant.Description,
		restaurant.MainImage,
		restaurant.IsVerified,
	).Scan(&restaurant.ID)
	if err != nil {
		return fmt.Errorf("insert restaurant: %w", err)
	}

	// Insert location
	query = `
		INSERT INTO restaurant_locations (restaurant_id, city, district, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5)`
	
	_, err = tx.ExecContext(ctx, query,
		restaurant.ID,
		restaurant.Location.Address.City,
		restaurant.Location.Address.District,
		restaurant.Location.Latitude,
		restaurant.Location.Longitude,
	)
	if err != nil {
		return fmt.Errorf("insert location: %w", err)
	}

	// Insert details
	query = `
		INSERT INTO restaurant_details (restaurant_id, rating, capacity, opening_hours)
		VALUES ($1, $2, $3, $4)`
	
	_, err = tx.ExecContext(ctx, query,
		restaurant.ID,
		restaurant.Details.Rating,
		restaurant.Details.Capacity,
		restaurant.Details.OpeningHours,
	)
	if err != nil {
		return fmt.Errorf("insert details: %w", err)
	}

	return tx.Commit()
}

func (r *RestaurantRepository) GetByID(ctx context.Context, id uint) (*domain.Restaurant, error) {
	query := `
		SELECT r.*, 
			   l.city, l.district, l.latitude, l.longitude,
			   d.rating, d.capacity, d.opening_hours
		FROM restaurants r
		LEFT JOIN restaurant_locations l ON l.restaurant_id = r.restaurant_id
		LEFT JOIN restaurant_details d ON d.restaurant_id = r.restaurant_id
		WHERE r.restaurant_id = $1`

	var restaurant domain.Restaurant
	err := r.db.GetContext(ctx, &restaurant, query, id)
	if err != nil {
		return nil, fmt.Errorf("get restaurant: %w", err)
	}

	return &restaurant, nil
}

func (r *RestaurantRepository) List(ctx context.Context, params dto.ListParams) ([]domain.Restaurant, int, error) {
	var total int
	query := `SELECT COUNT(*) FROM restaurants`
	err := r.db.GetContext(ctx, &total, query)
	if err != nil {
		return nil, 0, fmt.Errorf("count restaurants: %w", err)
	}

	query = `
		SELECT r.*, 
			   l.city, l.district, l.latitude, l.longitude,
			   d.rating, d.capacity, d.opening_hours
		FROM restaurants r
		LEFT JOIN restaurant_locations l ON l.restaurant_id = r.restaurant_id
		LEFT JOIN restaurant_details d ON d.restaurant_id = r.restaurant_id
		LIMIT $1 OFFSET $2`
	
	var restaurants []domain.Restaurant
	err = r.db.SelectContext(ctx, &restaurants, query, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list restaurants: %w", err)
	}

	return restaurants, total, nil
}

func (r *RestaurantRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM restaurants WHERE restaurant_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete restaurant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrRestaurantNotFound
	}

	return nil
}

func (r *RestaurantRepository) Update(ctx context.Context, restaurant *domain.Restaurant) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update restaurant
	query := `
		UPDATE restaurants 
		SET name = $1, description = $2, main_image = $3, is_verified = $4
		WHERE restaurant_id = $5`
	result, err := tx.ExecContext(ctx, query,
		restaurant.Name,
		restaurant.Description,
		restaurant.MainImage,
		restaurant.IsVerified,
		restaurant.ID,
	)
	if err != nil {
		return fmt.Errorf("update restaurant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrRestaurantNotFound
	}

	// Update location
	query = `
		UPDATE restaurant_locations 
		SET city = $1, district = $2, latitude = $3, longitude = $4
		WHERE restaurant_id = $5`
	_, err = tx.ExecContext(ctx, query,
		restaurant.Location.Address.City,
		restaurant.Location.Address.District,
		restaurant.Location.Latitude,
		restaurant.Location.Longitude,
		restaurant.ID,
	)
	if err != nil {
		return fmt.Errorf("update location: %w", err)
	}

	// Update details
	query = `
		UPDATE restaurant_details 
		SET rating = $1, capacity = $2, opening_hours = $3
		WHERE restaurant_id = $4`
	_, err = tx.ExecContext(ctx, query,
		restaurant.Details.Rating,
		restaurant.Details.Capacity,
		restaurant.Details.OpeningHours,
		restaurant.ID,
	)
	if err != nil {
		return fmt.Errorf("update details: %w", err)
	}

	return tx.Commit()
}

func (r *RestaurantRepository) ListByOwnerID(ctx context.Context, ownerID uint, params dto.ListParams) ([]domain.Restaurant, int, error) {
	var total int
	query := `SELECT COUNT(*) FROM restaurants WHERE owner_id = $1`
	err := r.db.GetContext(ctx, &total, query, ownerID)
	if err != nil {
		return nil, 0, fmt.Errorf("count restaurants: %w", err)
	}

	query = `
		SELECT r.*, 
			   l.city, l.district, l.latitude, l.longitude,
			   d.rating, d.capacity, d.opening_hours
		FROM restaurants r
		LEFT JOIN restaurant_locations l ON l.restaurant_id = r.restaurant_id
		LEFT JOIN restaurant_details d ON d.restaurant_id = r.restaurant_id
		WHERE r.owner_id = $1
		LIMIT $2 OFFSET $3`
	
	var restaurants []domain.Restaurant
	err = r.db.SelectContext(ctx, &restaurants, query, ownerID, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list restaurants: %w", err)
	}

	return restaurants, total, nil
} 