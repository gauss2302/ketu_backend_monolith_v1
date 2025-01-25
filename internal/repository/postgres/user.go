package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/repository/interfaces"
	"log"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) interfaces.UserRepository {
	if db == nil {
		panic("nil db provided to NewUserRepository")
	}
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	// Add logging
	log.Printf("Attempting to create user with email: %s", user.Email)

	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&exists)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return err
	}
	if exists {
		log.Printf("Email already exists: %s", user.Email)
		return domain.ErrEmailExists
	}

	query := `
	INSERT INTO users (username, email, password, name, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`

	err = r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	log.Printf("=== Getting user by email: %s ===", email)

	var user domain.User
	query := `SELECT id, username, email, password, name, created_at, updated_at FROM users WHERE email = $1`

	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrInvalidCredentials
		}
		log.Printf("Error querying user: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	query := `SELECT id, username, email, created_at, updated_at FROM users`

	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	query := `
        UPDATE users 
        SET username = $1, email = $2, updated_at = $3
        WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
