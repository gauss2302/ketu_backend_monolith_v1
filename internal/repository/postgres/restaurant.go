package postgres

import (
	"context"
	"fmt"
	"ketu_backend_monolith_v1/internal/domain"

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

func (r *RestaurantRepository) List(ctx context.Context) ([]domain.Restaurant, error) {
	query := `
		SELECT r.*, 
			   l.city, l.district, l.latitude, l.longitude,
			   d.rating, d.capacity, d.opening_hours
		FROM restaurants r
		LEFT JOIN restaurant_locations l ON l.restaurant_id = r.restaurant_id
		LEFT JOIN restaurant_details d ON d.restaurant_id = r.restaurant_id`

	var restaurants []domain.Restaurant
	err := r.db.SelectContext(ctx, &restaurants, query)
	if err != nil {
		return nil, fmt.Errorf("list restaurants: %w", err)
	}

	return restaurants, nil
} 