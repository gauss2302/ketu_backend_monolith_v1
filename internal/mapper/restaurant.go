package mapper

import (
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
)

func ToRestaurantDomain(req *dto.CreateRestaurantDTO) *domain.Restaurant {
	return &domain.Restaurant{
		OwnerID:     req.OwnerID,
		Name:        req.Name,
		Description: req.Description,
		MainImage:   req.MainImage,
		Images:      req.Images,
		Location: domain.RestaurantLocation{
			Address: domain.Address{
				City:     req.Location.Address.City,
				District: req.Location.Address.District,
			},
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
		Details: domain.RestaurantDetails{
			Capacity:     req.Details.Capacity,
			OpeningHours: req.Details.OpeningHours,
		},
	}
}

func ToRestaurantResponse(restaurant *domain.Restaurant) *dto.RestaurantResponse {
	return &dto.RestaurantResponse{
		ID:          restaurant.ID,
		Name:        restaurant.Name,
		Description: restaurant.Description,
		MainImage:   restaurant.MainImage,
		Images:      restaurant.Images,
		IsVerified:  restaurant.IsVerified,
		Location: dto.LocationResponse{
			Address: dto.AddressResponse{
				City:     restaurant.Location.Address.City,
				District: restaurant.Location.Address.District,
			},
			Latitude:  restaurant.Location.Latitude,
			Longitude: restaurant.Location.Longitude,
		},
		Details: dto.RestaurantDetailsResponse{
			Rating:       restaurant.Details.Rating,
			Capacity:     restaurant.Details.Capacity,
			OpeningHours: restaurant.Details.OpeningHours,
		},
		CreatedAt: restaurant.CreatedAt,
		UpdatedAt: restaurant.UpdatedAt,
	}
}

func ToRestaurantListResponse(restaurants []domain.Restaurant) []*dto.RestaurantResponse {
	var response []*dto.RestaurantResponse
	for _, restaurant := range restaurants {
		response = append(response, ToRestaurantResponse(&restaurant))
	}
	return response
} 