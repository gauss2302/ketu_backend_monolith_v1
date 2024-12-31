// service/place.go
package service

import (
	"context"
	"errors"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/mapper"
	repository "ketu_backend_monolith_v1/internal/repository/interfaces"
)

type PlaceService struct {
	repo repository.PlaceRepository
}

func NewPlaceService(repo repository.PlaceRepository) *PlaceService {
	return &PlaceService{
		repo: repo,
	}
}

// CreatePlace handles the creation of a new place
func (s *PlaceService) CreatePlace(ctx context.Context, createDTO *dto.PlaceCreateDTO) (*dto.PlaceResponseDTO, error) {
	if err := validateCreateDTO(createDTO); err != nil {
		return nil, err
	}

	place := mapper.ToPlaceDomain(createDTO)

	if err := s.repo.Create(ctx, place); err != nil {
		if errors.Is(err, domain.ErrPlaceExists) {
			return nil, domain.ErrPlaceExists
		}
		return nil, err
	}

	return mapper.ToPlaceResponseDTO(place), nil
}

// GetPlace retrieves a place by ID
func (s *PlaceService) GetPlace(ctx context.Context, id uint) (*dto.PlaceResponseDTO, error) {
	if id == 0 {
		return nil, domain.ErrInvalidID
	}

	place, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrPlaceNotFound) {
			return nil, domain.ErrPlaceNotFound
		}
		return nil, err
	}

	return mapper.ToPlaceResponseDTO(place), nil
}

// ListPlaces retrieves a paginated list of places
func (s *PlaceService) ListPlaces(ctx context.Context, params dto.PaginationParams) (*dto.PlaceListResponseDTO, error) {
	if err := validatePaginationParams(&params); err != nil {
		return nil, err
	}

	places, total, err := s.repo.List(ctx, params.Offset, params.Limit)
	if err != nil {
		return nil, err
	}

	return &dto.PlaceListResponseDTO{
		Places: mapper.ToPlaceListItemDTOs(places),
		Pagination: dto.PaginationResponse{
			Total:  total,
			Offset: params.Offset,
			Limit:  params.Limit,
		},
	}, nil
}

// SearchPlaces searches for places based on criteria
func (s *PlaceService) SearchPlaces(ctx context.Context, searchDTO *dto.PlaceSearchDTO) (*dto.PlaceListResponseDTO, error) {
	if err := validateSearchDTO(searchDTO); err != nil {
		return nil, err
	}

	criteria := dto.PlaceSearchCriteria{
		Name:      searchDTO.Name,
		City:      searchDTO.City,
		Province:  searchDTO.Province,
		Offset:    searchDTO.Offset,
		Limit:     searchDTO.Limit,
		SortBy:    searchDTO.SortBy,
		SortOrder: searchDTO.SortOrder,
	}

	places, total, err := s.repo.Search(ctx, criteria)
	if err != nil {
		return nil, err
	}

	return &dto.PlaceListResponseDTO{
		Places: mapper.ToPlaceListItemDTOs(places),
		Pagination: dto.PaginationResponse{
			Total:  total,
			Offset: searchDTO.Offset,
			Limit:  searchDTO.Limit,
		},
	}, nil
}

// UpdatePlace updates an existing place
func (s *PlaceService) UpdatePlace(ctx context.Context, id uint, updateDTO *dto.PlaceUpdateDTO) (*dto.PlaceResponseDTO, error) {
	if id == 0 {
		return nil, domain.ErrInvalidID
	}

	if err := validateUpdateDTO(updateDTO); err != nil {
		return nil, err
	}

	place, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrPlaceNotFound) {
			return nil, domain.ErrPlaceNotFound
		}
		return nil, err
	}

	mapper.UpdatePlaceDomain(place, updateDTO)

	if err := s.repo.Update(ctx, place); err != nil {
		return nil, err
	}

	return mapper.ToPlaceResponseDTO(place), nil
}

// DeletePlace removes a place by ID
func (s *PlaceService) DeletePlace(ctx context.Context, id uint) error {
	if id == 0 {
		return domain.ErrInvalidID
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrPlaceNotFound) {
			return domain.ErrPlaceNotFound
		}
		return err
	}

	return nil
}

// Validation helpers
func validateCreateDTO(dto *dto.PlaceCreateDTO) error {
	if dto == nil {
		return domain.ErrInvalidInput
	}
	if dto.Name == "" {
		return domain.ErrEmptyName
	}
	if dto.Description == "" {
		return domain.ErrEmptyDescription
	}
	if dto.Location.Address == "" {
		return domain.ErrEmptyAddress
	}
	return nil
}

func validateUpdateDTO(dto *dto.PlaceUpdateDTO) error {
	if dto == nil {
		return domain.ErrEmailExists
	}
	if dto.Name != nil && *dto.Name == "" {
		return domain.ErrEmptyName
	}
	if dto.Description != nil && *dto.Description == "" {
		return domain.ErrEmptyDescription
	}
	return nil
}

func validatePaginationParams(params *dto.PaginationParams) error {
	if params == nil {
		return domain.ErrInvalidInput
	}
	if params.Limit <= 0 {
		params.Limit = 10 // default limit
	}
	if params.Offset < 0 {
		params.Offset = 0
	}
	return nil
}

func validateSearchDTO(dto *dto.PlaceSearchDTO) error {
	if dto == nil {
		return domain.ErrInvalidInput
	}
	if dto.Limit <= 0 {
		dto.Limit = 10 // default limit
	}
	if dto.Offset < 0 {
		dto.Offset = 0
	}
	return nil
}
