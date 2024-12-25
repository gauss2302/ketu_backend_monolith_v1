package mapper

import (
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/handler/dto"

	"time"
)

// ToPlaceDomain converts a PlaceCreateDTO to a domain Place
func ToPlaceDomain(createDTO *dto.PlaceCreateDTO) *domain.Place {
	return &domain.Place{
		Name:        createDTO.Name,
		Description: createDTO.Description,
		Location: domain.Location{
			Address:  createDTO.Location.Address,
			City:     createDTO.Location.City,
			Province: createDTO.Location.Province,
		},
		MainImage: createDTO.MainImage,
		Images:    createDTO.Images,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// UpdatePlaceDomain updates a domain Place with data from PlaceUpdateDTO
func UpdatePlaceDomain(place *domain.Place, updateDTO *dto.PlaceUpdateDTO) {
	if updateDTO.Name != nil {
		place.Name = *updateDTO.Name
	}
	if updateDTO.Description != nil {
		place.Description = *updateDTO.Description
	}
	if updateDTO.Location != nil {
		place.Location = domain.Location{
			Address:  updateDTO.Location.Address,
			City:     updateDTO.Location.City,
			Province: updateDTO.Location.Province,
		}
	}
	if updateDTO.MainImage != nil {
		place.MainImage = *updateDTO.MainImage
	}
	if updateDTO.Images != nil {
		place.Images = updateDTO.Images
	}
	place.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
}

// ToPlaceResponseDTO converts a domain Place to PlaceResponseDTO
func ToPlaceResponseDTO(place *domain.Place) *dto.PlaceResponseDTO {
	return &dto.PlaceResponseDTO{
		ID:          place.ID,
		Name:        place.Name,
		Description: place.Description,
		Location: dto.LocationDTO{
			Address:  place.Location.Address,
			City:     place.Location.City,
			Province: place.Location.Province,
		},
		MainImage: place.MainImage,
		Images:    place.Images,
		CreatedAt: place.CreatedAt,
	}
}

// ToPlaceListItemDTO converts a domain Place to PlaceListItemDTO
func ToPlaceListItemDTO(place *domain.Place) *dto.PlaceListItemDTO {
	return &dto.PlaceListItemDTO{
		ID:        place.ID,
		Name:      place.Name,
		City:      place.Location.City,
		Province:  place.Location.Province,
		MainImage: place.MainImage,
	}
}

// ToPlaceListItemDTOs converts a slice of domain Places to PlaceListItemDTOs
func ToPlaceListItemDTOs(places []*domain.Place) []*dto.PlaceListItemDTO {
	dtos := make([]*dto.PlaceListItemDTO, len(places))
	for i, place := range places {
		dtos[i] = ToPlaceListItemDTO(place)
	}
	return dtos
}
