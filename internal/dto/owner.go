package dto

import "time"

type OwnerRegisterRequestDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type OwnerLoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type OwnerResponseDTO struct {
	ID        uint      `json:"owner_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type OwnerAuthResponseDTO struct {
	Owner       *OwnerResponseDTO `json:"owner"`
	AccessToken string            `json:"access_token"`
	ExpiresIn   int64            `json:"expires_in"`
}

type OwnerCreateDTO struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
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