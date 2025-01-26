package interfaces

import (
	"context"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRestaurantRepository is a mock implementation of RestaurantRepository
type MockRestaurantRepository struct {
	mock.Mock
}

func (m *MockRestaurantRepository) Create(ctx context.Context, restaurant *domain.Restaurant) error {
	args := m.Called(ctx, restaurant)
	return args.Error(0)
}

func (m *MockRestaurantRepository) GetByID(ctx context.Context, id uint) (*domain.Restaurant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) List(ctx context.Context, params dto.ListParams) ([]domain.Restaurant, int, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]domain.Restaurant), args.Int(1), args.Error(2)
}

func (m *MockRestaurantRepository) Update(ctx context.Context, restaurant *domain.Restaurant) error {
	args := m.Called(ctx, restaurant)
	return args.Error(0)
}

func (m *MockRestaurantRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRestaurantRepository) ListByOwnerID(ctx context.Context, ownerID uint, params dto.ListParams) ([]domain.Restaurant, int, error) {
	args := m.Called(ctx, ownerID, params)
	return args.Get(0).([]domain.Restaurant), args.Int(1), args.Error(2)
}

func TestRestaurantRepository(t *testing.T) {
	mockRepo := new(MockRestaurantRepository)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		restaurant := &domain.Restaurant{
			Name:        "Test Restaurant",
			Description: "Test Description",
			OwnerID:     1,
		}

		mockRepo.On("Create", ctx, restaurant).Return(nil)

		err := mockRepo.Create(ctx, restaurant)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByID", func(t *testing.T) {
		expectedRestaurant := &domain.Restaurant{
			ID:          1,
			Name:        "Test Restaurant",
			Description: "Test Description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.On("GetByID", ctx, uint(1)).Return(expectedRestaurant, nil)

		restaurant, err := mockRepo.GetByID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedRestaurant, restaurant)
		mockRepo.AssertExpectations(t)
	})

	t.Run("List", func(t *testing.T) {
		params := dto.ListParams{
			Offset: 0,
			Limit:  10,
		}

		expectedRestaurants := []domain.Restaurant{
			{
				ID:          1,
				Name:        "Restaurant 1",
				Description: "Description 1",
			},
			{
				ID:          2,
				Name:        "Restaurant 2",
				Description: "Description 2",
			},
		}

		mockRepo.On("List", ctx, params).Return(expectedRestaurants, 2, nil)

		restaurants, total, err := mockRepo.List(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, expectedRestaurants, restaurants)
		assert.Equal(t, 2, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update", func(t *testing.T) {
		restaurant := &domain.Restaurant{
			ID:          1,
			Name:        "Updated Restaurant",
			Description: "Updated Description",
		}

		mockRepo.On("Update", ctx, restaurant).Return(nil)

		err := mockRepo.Update(ctx, restaurant)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		mockRepo.On("Delete", ctx, uint(1)).Return(nil)

		err := mockRepo.Delete(ctx, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ListByOwnerID", func(t *testing.T) {
		params := dto.ListParams{
			Offset: 0,
			Limit:  10,
		}

		expectedRestaurants := []domain.Restaurant{
			{
				ID:      1,
				OwnerID: 1,
				Name:    "Owner's Restaurant 1",
			},
			{
				ID:      2,
				OwnerID: 1,
				Name:    "Owner's Restaurant 2",
			},
		}

		mockRepo.On("ListByOwnerID", ctx, uint(1), params).Return(expectedRestaurants, 2, nil)

		restaurants, total, err := mockRepo.ListByOwnerID(ctx, 1, params)
		assert.NoError(t, err)
		assert.Equal(t, expectedRestaurants, restaurants)
		assert.Equal(t, 2, total)
		mockRepo.AssertExpectations(t)
	})
} 