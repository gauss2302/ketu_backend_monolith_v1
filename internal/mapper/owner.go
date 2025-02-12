package mapper

import (
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
)

func ToOwnerResponseDTO(owner *domain.Owner) *dto.OwnerResponseDTO {
	return &dto.OwnerResponseDTO{
		ID:        owner.ID,
		Name:      owner.Name,
		Email:     owner.Email,
		Phone:     owner.Phone,
		CreatedAt: owner.CreatedAt,
	}
} 