package dto

type OwnerResponseDTO struct {
	ID          uint                  `json:"owner_id"`
	Name        string                `json:"name"`
	Email       string                `json:"email"`
	Phone       string                `json:"phone"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
	Restaurants []RestaurantBasicDTO  `json:"restaurants,omitempty"`
}

type OwnerCreateDTO struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
}

type OwnerUpdateDTO struct {
	Name  string `json:"name,omitempty" validate:"omitempty"`
	Email string `json:"email,omitempty" validate:"omitempty,email"`
	Phone string `json:"phone,omitempty" validate:"omitempty"`
}

type RestaurantBasicDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}