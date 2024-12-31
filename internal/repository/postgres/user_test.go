package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"ketu_backend_monolith_v1/internal/domain"
	repository "ketu_backend_monolith_v1/internal/repository/interfaces"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (sqlmock.Sqlmock, repository.UserRepository) {
    mockDB, mock, err := sqlmock.New()
    require.NoError(t, err, "Failed to create mock database")

    sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
    repo := NewUserRepository(sqlxDB)
    
    return mock, repo
}

// Helper function to create test user
func createTestUser() *domain.User {
    now := time.Now().UTC()
    return &domain.User{
        Username:  "testuser",
        Email:     "test@example.com",
        Password:  "hashedpassword",
        Name:      "Test User",
        CreatedAt: now,
        UpdatedAt: now,
    }
}

func TestCreate(t *testing.T) {
	mock, repo := setupTestDB(t)
	ctx := context.Background()
	user := createTestUser()

	tests := []struct {
		 name    string
		 user    *domain.User
		 setup   func(mock sqlmock.Sqlmock)
		 wantErr error
	}{
		 {
			  name: "success",
			  user: user,
			  setup: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT EXISTS").
						 WithArgs(user.Email).
						 WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

					mock.ExpectQuery("INSERT INTO users").
						 WithArgs(user.Username, user.Email, user.Password, user.Name, user.CreatedAt, user.UpdatedAt).
						 WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			  },
			  wantErr: nil,
		 },
		 {
			  name: "email_exists",
			  user: user,
			  setup: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT EXISTS").
						 WithArgs(user.Email).
						 WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			  },
			  wantErr: domain.ErrEmailExists,
		 },
		 {
			  name: "database_error",
			  user: user,
			  setup: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT EXISTS").
						 WithArgs(user.Email).
						 WillReturnError(sql.ErrConnDone)
			  },
			  wantErr: sql.ErrConnDone,
		 },
	}

	for _, tt := range tests {
		 t.Run(tt.name, func(t *testing.T) {
			  tt.setup(mock)
			  
			  err := repo.Create(ctx, tt.user)
			  
			  if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
			  } else {
					assert.NoError(t, err)
					assert.Equal(t, uint(1), tt.user.ID)
			  }
		 })
	}
}

func TestGetByEmail(t *testing.T) {
	mock, repo := setupTestDB(t)
	ctx := context.Background()
	user := createTestUser()

	tests := []struct {
		 name    string
		 email   string
		 setup   func(mock sqlmock.Sqlmock)
		 want    *domain.User
		 wantErr error
	}{
		 {
			  name:  "success",
			  email: user.Email,
			  setup: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{
						 "id", "username", "email", "password", "name", "created_at", "updated_at",
					}).AddRow(
						 1, user.Username, user.Email, user.Password, user.Name,
						 user.CreatedAt, user.UpdatedAt,
					)

					mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
						 WithArgs(user.Email).
						 WillReturnRows(rows)
			  },
			  want:    user,
			  wantErr: nil,
		 },
		 {
			  name:  "user_not_found",
			  email: "nonexistent@example.com",
			  setup: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
						 WithArgs("nonexistent@example.com").
						 WillReturnError(sql.ErrNoRows)
			  },
			  want:    nil,
			  wantErr: domain.ErrInvalidCredentials,
		 },
	}

	for _, tt := range tests {
		 t.Run(tt.name, func(t *testing.T) {
			  tt.setup(mock)

			  got, err := repo.GetByEmail(ctx, tt.email)
			  
			  if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
					assert.Nil(t, got)
			  } else {
					assert.NoError(t, err)
					assert.NotNil(t, got)
					assert.Equal(t, tt.want.Email, got.Email)
					assert.Equal(t, tt.want.Username, got.Username)
			  }
		 })
	}
}

func TestUpdate(t *testing.T) {
	mock, repo := setupTestDB(t)
	ctx := context.Background()
	user := createTestUser()
	user.ID = 1

	tests := []struct {
		 name    string
		 user    *domain.User
		 setup   func(mock sqlmock.Sqlmock)
		 wantErr error
	}{
		 {
			  name: "success",
			  user: user,
			  setup: func(mock sqlmock.Sqlmock) {
					mock.ExpectExec("UPDATE users").
						 WithArgs(
							  user.Username,
							  user.Email,
							  user.UpdatedAt,
							  user.ID,
						 ).
						 WillReturnResult(sqlmock.NewResult(1, 1))
			  },
			  wantErr: nil,
		 },
		 {
			  name: "user_not_found",
			  user: user,
			  setup: func(mock sqlmock.Sqlmock) {
					mock.ExpectExec("UPDATE users").
						 WithArgs(
							  user.Username,
							  user.Email,
							  user.UpdatedAt,
							  user.ID,
						 ).
						 WillReturnResult(sqlmock.NewResult(0, 0))
			  },
			  wantErr: domain.ErrUserNotFound,
		 },
	}

	for _, tt := range tests {
		 t.Run(tt.name, func(t *testing.T) {
			  tt.setup(mock)

			  err := repo.Update(ctx, tt.user)

			  if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
			  } else {
					assert.NoError(t, err)
			  }
		 })
	}
}