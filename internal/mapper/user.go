// internal/mapper/user.go
package mapper

import (
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"time"
)

func ToUserDomain(req *dto.UserCreateDTO) *domain.User {
	now := time.Now()
	return &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func UpdateUserDomain(user *domain.User, req *dto.UserUpdateDTO) {
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	user.UpdatedAt = time.Now()
}

func ToUserResponseDTO(user *domain.User) *dto.UserResponseDTO {
	return &dto.UserResponseDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func ToUserListResponseDTO(users []domain.User) []*dto.UserResponseDTO {
	dtos := make([]*dto.UserResponseDTO, len(users))
	for i, user := range users {
		dtos[i] = ToUserResponseDTO(&user)
	}
	return dtos
}
