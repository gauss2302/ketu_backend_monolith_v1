package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"ketu_backend_monolith_v1/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *UserRepo) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewUserRepository(sqlxDB)
	return sqlxDB, mock, repo
}

func TestCreate(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	ctx := context.Background()

	user := &domain.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(user.Email).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.Create(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), user.ID)
	})

	t.Run("email exists", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(user.Email).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		err := repo.Create(ctx, user)
		assert.Equal(t, domain.ErrEmailExists, err)
	})
}

func TestGetByID(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow(1, "testuser", "test@example.com", time.Now(), time.Now())

		mock.ExpectQuery("SELECT (.+) FROM users WHERE").
			WithArgs(uint(1)).
			WillReturnRows(rows)

		user, err := repo.GetByID(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE").
			WithArgs(uint(999)).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByID(ctx, 999)
		assert.Equal(t, domain.ErrUserNotFound, err)
		assert.Nil(t, user)
	})
}

func TestGetByEmail(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
			AddRow(1, "testuser", "test@example.com", "hashedpassword", time.Now(), time.Now())

		mock.ExpectQuery("SELECT (.+) FROM users WHERE").
			WithArgs("test@example.com").
			WillReturnRows(rows)

		user, err := repo.GetByEmail(ctx, "test@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
	})
}

func TestGetAll(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
		AddRow(1, "user1", "user1@example.com", time.Now(), time.Now()).
		AddRow(2, "user2", "user2@example.com", time.Now(), time.Now())

	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnRows(rows)

	users, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUpdate(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	ctx := context.Background()

	user := &domain.User{
		ID:        1,
		Username:  "updated",
		Email:     "updated@example.com",
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE users").
			WithArgs(user.Username, user.Email, user.UpdatedAt, user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(ctx, user)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE users").
			WithArgs(user.Username, user.Email, user.UpdatedAt, user.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(ctx, user)
		assert.Equal(t, domain.ErrUserNotFound, err)
	})
}

func TestDelete(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM users").
			WithArgs(uint(1)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(ctx, 1)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM users").
			WithArgs(uint(999)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(ctx, 999)
		assert.Equal(t, domain.ErrUserNotFound, err)
	})
}
