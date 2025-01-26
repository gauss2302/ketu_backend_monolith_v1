package dto

import "time"

// Request DTOs
type CreateRestaurantDTO struct {
	OwnerID     uint     `json:"owner_id" validate:"required"`
	Name        string   `json:"name" validate:"required,min=2,max=100"`
	Description string   `json:"description" validate:"required,min=10,max=1000"`
	MainImage   string   `json:"main_image" validate:"required,url"`
	Images      []string `json:"images" validate:"dive,url"`
	Location    RestaurantLocationDTO `json:"location" validate:"required"`
	Details     RestaurantDetailsDTO `json:"details" validate:"required"`
}

type RestaurantLocationDTO struct {
	Address   AddressDTO `json:"address" validate:"required"`
	Latitude  float64    `json:"latitude" validate:"required,latitude"`
	Longitude float64    `json:"longitude" validate:"required,longitude"`
}

type AddressDTO struct {
	City     string `json:"city" validate:"required,min=2,max=50"`
	District string `json:"district" validate:"required,min=2,max=50"`
}

type RestaurantDetailsDTO struct {
	Capacity     int     `json:"capacity" validate:"required,min=1"`
	OpeningHours string  `json:"opening_hours" validate:"required"`
}

// Response DTOs
type RestaurantResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	MainImage   string                  `json:"main_image"`
	Images      []string                `json:"images"`
	IsVerified  bool                    `json:"is_verified"`
	Location    LocationResponse        `json:"location"`
	Details     RestaurantDetailsResponse `json:"details"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

type LocationResponse struct {
	Address    AddressResponse `json:"address"`
	Latitude   float64         `json:"latitude"`
	Longitude  float64         `json:"longitude"`
}

type AddressResponse struct {
	City     string `json:"city"`
	District string `json:"district"`
}

type RestaurantDetailsResponse struct {
	Rating       float64 `json:"rating"`
	Capacity     int     `json:"capacity"`
	OpeningHours string  `json:"opening_hours"`
}

// Update DTO
type UpdateRestaurantDTO struct {
	Name        *string   `json:"name" validate:"omitempty,min=2,max=100"`
	Description *string   `json:"description" validate:"omitempty,min=10,max=1000"`
	MainImage   *string   `json:"main_image" validate:"omitempty,url"`
	Images      []string  `json:"images" validate:"omitempty,dive,url"`
	Location    *RestaurantLocationDTO `json:"location" validate:"omitempty"`
	Details     *RestaurantDetailsDTO `json:"details" validate:"omitempty"`
}

