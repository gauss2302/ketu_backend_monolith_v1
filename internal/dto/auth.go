package dto

// Users
type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequestDTO struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthResponseDTO struct {
	User         *UserResponseDTO `json:"user"`
	AccessToken  string           `json:"accessToken"`
	ExpiresIn    int64            `json:"expiresIn"`
}

//Tokens
type TokenRefreshResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}
