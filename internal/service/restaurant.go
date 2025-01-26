// service/restaurant.go
package service

import (
	"context"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/mapper"
	repository "ketu_backend_monolith_v1/internal/repository/interfaces"
)

type RestaurantService struct {
	repo repository.RestaurantRepository
}

func NewRestaurantService(repo repository.RestaurantRepository) *RestaurantService {
	return &RestaurantService{
		repo: repo,
	}
}

func (s *RestaurantService) Create(ctx context.Context, createDTO *dto.CreateRestaurantDTO) (*dto.RestaurantResponse, error) {
	if err := validateCreateRestaurantDTO(createDTO); err != nil {
		return nil, err
	}

	restaurant := mapper.ToRestaurantDomain(createDTO)

	if err := s.repo.Create(ctx, restaurant); err != nil {
		return nil, err
	}

	return mapper.ToRestaurantResponse(restaurant), nil
}

func (s *RestaurantService) GetByID(ctx context.Context, id uint) (*dto.RestaurantResponse, error) {
	restaurant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.ToRestaurantResponse(restaurant), nil
}

func (s *RestaurantService) List(ctx context.Context, params dto.ListParams) ([]*dto.RestaurantResponse, int, error) {
	restaurants, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	return mapper.ToRestaurantListResponse(restaurants), total, nil
}

func (s *RestaurantService) ListByOwnerID(ctx context.Context, ownerID uint, params dto.ListParams) ([]*dto.RestaurantResponse, int, error) {
	restaurants, total, err := s.repo.ListByOwnerID(ctx, ownerID, params)
	if err != nil {
		return nil, 0, err
	}

	return mapper.ToRestaurantListResponse(restaurants), total, nil
}

func (s *RestaurantService) Update(ctx context.Context, id uint, ownerID uint, updateDTO *dto.UpdateRestaurantDTO) (*dto.RestaurantResponse, error) {
	restaurant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if restaurant.OwnerID != ownerID {
		return nil, domain.ErrUnauthorized
	}

	if err := validateUpdateRestaurantDTO(updateDTO); err != nil {
		return nil, err
	}

	// Update fields
	if updateDTO.Name != nil {
		restaurant.Name = *updateDTO.Name
	}
	if updateDTO.Description != nil {
		restaurant.Description = *updateDTO.Description
	}
	if updateDTO.MainImage != nil {
		restaurant.MainImage = *updateDTO.MainImage
	}
	if updateDTO.Location != nil {
		restaurant.Location.Address = domain.Address{
			City:     updateDTO.Location.Address.City,
			District: updateDTO.Location.Address.District,
		}
		restaurant.Location.Latitude = updateDTO.Location.Latitude
		restaurant.Location.Longitude = updateDTO.Location.Longitude
	}
	if updateDTO.Details != nil {
		restaurant.Details.Capacity = updateDTO.Details.Capacity
		restaurant.Details.OpeningHours = updateDTO.Details.OpeningHours
	}

	if err := s.repo.Update(ctx, restaurant); err != nil {
		return nil, err
	}

	return mapper.ToRestaurantResponse(restaurant), nil
}

func (s *RestaurantService) Delete(ctx context.Context, id uint, ownerID uint) error {
	restaurant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if restaurant.OwnerID != ownerID {
		return domain.ErrUnauthorized
	}

	return s.repo.Delete(ctx, id)
}

// Validation helpers
func validateCreateRestaurantDTO(dto *dto.CreateRestaurantDTO) error {
	if dto == nil {
		return domain.ErrInvalidInput
	}
	if dto.Name == "" {
		return domain.ErrEmptyName
	}
	if dto.Description == "" {
		return domain.ErrEmptyDescription
	}
	return nil
}

func validateUpdateRestaurantDTO(dto *dto.UpdateRestaurantDTO) error {
	if dto == nil {
		return domain.ErrInvalidInput
	}
	if dto.Name != nil && *dto.Name == "" {
		return domain.ErrEmptyName
	}
	if dto.Description != nil && *dto.Description == "" {
		return domain.ErrEmptyDescription
	}
	return nil
}

