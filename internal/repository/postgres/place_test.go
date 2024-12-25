package postgres

// import (
// 	"context"
// 	"database/sql"
// 	"ketu_backend_monolith_v1/internal/domain"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestPlaceRepo_CreatePlace(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	repo := NewPlaceRepository(db)

// 	tests := []struct {
// 		name    string
// 		place   *domain.Place
// 		mock    func()
// 		wantErr bool
// 	}{
// 		{
// 			name: "Success",
// 			place: &domain.Place{
// 				Name:        "Test Place",
// 				Description: "Test Description",
// 				Location: domain.Location{
// 					Address:  "Test Address",
// 					City:     "Test City",
// 					Province: "Test Province",
// 				},
// 				MainImage: "test.jpg",
// 				CreatedAt: time.Now().Format(time.RFC3339),
// 				UpdatedAt: time.Now().Format(time.RFC3339),
// 			},
// 			mock: func() {
// 				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
// 				mock.ExpectQuery("INSERT INTO places").
// 					WithArgs("Test Place", "Test Description", sqlmock.AnyArg(), "test.jpg", sqlmock.AnyArg(), sqlmock.AnyArg()).
// 					WillReturnRows(rows)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Database Error",
// 			place: &domain.Place{
// 				Name:        "Test Place",
// 				Description: "Test Description",
// 				Location: domain.Location{
// 					Address:  "Test Address",
// 					City:     "Test City",
// 					Province: "Test Province",
// 				},
// 			},
// 			mock: func() {
// 				mock.ExpectQuery("INSERT INTO places").
// 					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
// 					WillReturnError(sql.ErrConnDone)
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()

// 			err := repo.CreatePlace(context.Background(), tt.place)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.NotZero(t, tt.place.ID)
// 			}
// 		})
// 	}
// }

// func TestPlaceRepo_GetPlaceByID(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	repo := NewPlaceRepository(db)

// 	tests := []struct {
// 		name    string
// 		id      uint
// 		mock    func()
// 		want    *domain.Place
// 		wantErr bool
// 	}{
// 		{
// 			name: "Success",
// 			id:   1,
// 			mock: func() {
// 				rows := sqlmock.NewRows([]string{"id", "name", "description", "location", "main_image", "created_at", "updated_at"}).
// 					AddRow(1, "Test Place", "Test Description", sqlmock.AnyArg(), "test.jpg", time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339))
// 				mock.ExpectQuery("SELECT (.+) FROM places").
// 					WithArgs(1).
// 					WillReturnRows(rows)
// 			},
// 			want: &domain.Place{
// 				ID:          1,
// 				Name:        "Test Place",
// 				Description: "Test Description",
// 				MainImage:   "test.jpg",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Not Found",
// 			id:   999,
// 			mock: func() {
// 				mock.ExpectQuery("SELECT (.+) FROM places").
// 					WithArgs(999).
// 					WillReturnError(sql.ErrNoRows)
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()

// 			got, err := repo.GetPlaceByID(context.Background(), tt.id)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				return
// 			}

// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.want.ID, got.ID)
// 			assert.Equal(t, tt.want.Name, got.Name)
// 			assert.Equal(t, tt.want.Description, got.Description)
// 			assert.Equal(t, tt.want.MainImage, got.MainImage)
// 		})
// 	}
// }

// func TestPlaceRepo_UpdatePlaceByID(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	repo := NewPlaceRepository(db)

// 	tests := []struct {
// 		name    string
// 		place   *domain.Place
// 		mock    func()
// 		wantErr bool
// 	}{
// 		{
// 			name: "Success",
// 			place: &domain.Place{
// 				ID:          1,
// 				Name:        "Updated Place",
// 				Description: "Updated Description",
// 				Location: domain.Location{
// 					Address:  "Updated Address",
// 					City:     "Updated City",
// 					Province: "Updated Province",
// 				},
// 				MainImage: "updated.jpg",
// 				UpdatedAt: time.Now().Format(time.RFC3339),
// 			},
// 			mock: func() {
// 				mock.ExpectExec("UPDATE places").
// 					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
// 					WillReturnResult(sqlmock.NewResult(1, 1))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Not Found",
// 			place: &domain.Place{
// 				ID: 999,
// 			},
// 			mock: func() {
// 				mock.ExpectExec("UPDATE places").
// 					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 999).
// 					WillReturnResult(sqlmock.NewResult(0, 0))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()

// 			err := repo.UpdatePlaceByID(context.Background(), tt.place)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestPlaceRepo_DeletePlaceByID(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	repo := NewPlaceRepository(db)

// 	tests := []struct {
// 		name    string
// 		id      uint
// 		mock    func()
// 		wantErr bool
// 	}{
// 		{
// 			name: "Success",
// 			id:   1,
// 			mock: func() {
// 				mock.ExpectExec("DELETE FROM places").
// 					WithArgs(1).
// 					WillReturnResult(sqlmock.NewResult(1, 1))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Not Found",
// 			id:   999,
// 			mock: func() {
// 				mock.ExpectExec("DELETE FROM places").
// 					WithArgs(999).
// 					WillReturnResult(sqlmock.NewResult(0, 0))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()

// 			err := repo.DeletePlaceByID(context.Background(), tt.id)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }
