package dto

type LocationDTO struct {
	Address  string `json:"address"`
	City     string `json:"city"`
	Province string `json:"province"`
}

// PlaceCreateDTO represents the data needed to create a new place
type PlaceCreateDTO struct {
	Name        string      `json:"name" validate:"required,min=3,max=100"`
	Description string      `json:"description" validate:"required,min=10,max=1000"`
	Location    LocationDTO `json:"location" validate:"required"`
	MainImage   string      `json:"main_image" validate:"required,url"`
	Images      []string    `json:"images" validate:"dive,url"`
}

// PlaceUpdateDTO represents the data needed to update an existing place
type PlaceUpdateDTO struct {
	Name        *string      `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string      `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Location    *LocationDTO `json:"location,omitempty"`
	MainImage   *string      `json:"main_image,omitempty" validate:"omitempty,url"`
	Images      []string     `json:"images,omitempty" validate:"omitempty,dive,url"`
}

// PlaceResponseDTO represents the data returned to clients
type PlaceResponseDTO struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Location    LocationDTO `json:"location"`
	MainImage   string      `json:"main_image"`
	Images      []string    `json:"images,omitempty"`
	CreatedAt   string      `json:"created_at"`
}

// PlaceListItemDTO represents a condensed version of place for list views
type PlaceListItemDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	City      string `json:"city"`
	Province  string `json:"province"`
	MainImage string `json:"main_image"`
}

type PaginationParams struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=1,max=100"`
}

type PaginationResponse struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type PlaceSearchDTO struct {
	PaginationParams
	Name      string `json:"name,omitempty"`
	City      string `json:"city,omitempty"`
	Province  string `json:"province,omitempty"`
	SortBy    string `json:"sort_by,omitempty" validate:"omitempty,oneof=name created_at updated_at"`
	SortOrder string `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
}

type PlaceListResponseDTO struct {
	Places     []*PlaceListItemDTO `json:"places"`
	Pagination PaginationResponse  `json:"pagination"`
}

type PlaceSearchCriteria struct {
	Name      string
	City      string
	Province  string
	Offset    int
	Limit     int
	SortBy    string
	SortOrder string
}
