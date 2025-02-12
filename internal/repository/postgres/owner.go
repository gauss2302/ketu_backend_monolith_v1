package postgres

import (
	"context"
	"fmt"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/pkg/database"
	"ketu_backend_monolith_v1/internal/repository/interfaces"
	"log"
	"time"
)

type OwnerRepo struct {
	db *database.DB
}

func NewOwnerRepository(db *database.DB) interfaces.OwnerRepository {
	if db == nil {
		panic("nil db provided to NewOwnerRepository")
	}
	return &OwnerRepo{
		db: db,
	}
}

func (r *OwnerRepo) Create(ctx context.Context, owner *domain.Owner) error {
	log.Printf("Attempting to create an owner with email: %s", owner.Email)

	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM owners WHERE email = $1)", owner.Email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("checking email existence: %w", err)
	}

	if exists {
		return domain.ErrEmailExists
	}

	log.Printf("Debug - Password value before insert: %s", owner.Password)

	query := `
	INSERT INTO owners (name, email, phone, password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING owner_id`

	now := time.Now()
	owner.CreatedAt = now
	owner.UpdatedAt = now

	err = r.db.QueryRowContext(ctx, query,
		owner.Name,
		owner.Email,
		owner.Phone,
		owner.Password,
		owner.CreatedAt,
		owner.UpdatedAt,
	).Scan(&owner.ID)

	if err != nil {
		return fmt.Errorf("creating owner: %w", err)
	}

	return nil
}

func (r *OwnerRepo) GetByID(ctx context.Context, id uint) (*domain.Owner, error) {
	query := `
	SELECT o.*, 
		   COALESCE(json_agg(r.*) FILTER (WHERE r.restaurant_id IS NOT NULL), '[]') as restaurants
	FROM owners o
	LEFT JOIN restaurants r ON r.owner_id = o.owner_id
	WHERE o.owner_id = $1
	GROUP BY o.owner_id`

	var owner domain.Owner
	err := r.db.GetContext(ctx, &owner, query, id)
	if err != nil {
		return nil, fmt.Errorf("get owner: %w", err)
	}

	return &owner, nil
}

func (r *OwnerRepo) GetByEmail(ctx context.Context, email string) (*domain.Owner, error) {
	query := `
	SELECT owner_id, name, email, phone, password, created_at, updated_at
	FROM owners 
	WHERE email = $1`

	var owner domain.Owner
	err := r.db.GetContext(ctx, &owner, query, email)
	if err != nil {
		return nil, fmt.Errorf("get owner by email: %w", err)
	}

	return &owner, nil
}

func (r *OwnerRepo) GetAll(ctx context.Context) ([]domain.Owner, error) {
	query := `
	SELECT o.*, 
		   COALESCE(json_agg(r.*) FILTER (WHERE r.restaurant_id IS NOT NULL), '[]') as restaurants
	FROM owners o
	LEFT JOIN restaurants r ON r.owner_id = o.owner_id
	GROUP BY o.owner_id`

	var owners []domain.Owner
	err := r.db.SelectContext(ctx, &owners, query)
	if err != nil {
		return nil, fmt.Errorf("list owners: %w", err)
	}

	return owners, nil
}

func (r *OwnerRepo) Update(ctx context.Context, owner *domain.Owner) error {
	query := `
	UPDATE owners 
	SET name = $1, 
		email = $2, 
		phone = $3,
		updated_at = $4
	WHERE owner_id = $5`

	owner.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		owner.Name,
		owner.Email,
		owner.Phone,
		owner.UpdatedAt,
		owner.ID,
	)
	if err != nil {
		return fmt.Errorf("update owner: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrOwnerNotFound
	}

	return nil
}

func (r *OwnerRepo) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM owners WHERE owner_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete owner: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrOwnerNotFound
	}

	return nil
}
