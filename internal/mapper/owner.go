package mapper

import (
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"time"
)


func ToOwnerDomain(req *dto.OwnerCreateDTO) *domain.Owner {
	now := time.Now()
	return &domain.Owner{
		Name: req.Name,
		Email: req.Email,
		Phone: req.Phone,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func UpdateOwnerDomain(owner *domain.Owner, req *dto.OwnerUpdateDTO) {
	if req.Name != nil {
		owner.Name = *req.Name
	}
	if req.Email != nil {
		owner.Email = *req.Email
	}
	owner.UpdatedAt = time.Now()
}

func ToOwnerResponseDTO(owner *domain.Owner) *dto.OwnerResponseDTO {
	return &dto.OwnerResponseDTO{
		ID: owner.ID,
		Name: owner.Name,
		Email: owner.Email,
		Phone: owner.Phone,
		CreatedAt: owner.CreatedAt,
		UpdatedAt: owner.UpdatedAt,
	}
}

func ToOwnerListResponseDTO(owners []domain.Owner) []*dto.OwnerResponseDTO {
	dtos := make([]*dto.OwnerResponseDTO, len(owners))

	for i, owner := range owners {
		dtos[i] = ToOwnerResponseDTO(&owner)
	}
	return dtos
}
