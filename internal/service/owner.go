package service

import (
	"context"
	"errors"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/mapper"
	repository "ketu_backend_monolith_v1/internal/repository/interfaces"

	"golang.org/x/crypto/bcrypt"
)


type OwnerService struct {
	ownerRepo repository.OwnerRepository
}

func NewOwnerService(ownerRepo repository.OwnerRepository) *OwnerService {
	return &OwnerService {
		ownerRepo: ownerRepo,
	}
}

func (s *OwnerService) Create(ctx context.Context, req *dto.OwnerCreateDTO) (*dto.OwnerResponseDTO, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	owner := mapper.ToOwnerDomain(req)
	owner.Password = string(hashedPassword)

	if err := s.ownerRepo.Create(ctx, owner); err != nil {
		if errors.Is(err, domain.ErrEmptyName) {
			return nil, domain.ErrEmptyName
		}
		return nil, err
	}

	return mapper.ToOwnerResponseDTO(owner), nil

}

func (s *OwnerService) GetByID(ctx context.Context, id uint) (*dto.OwnerResponseDTO, error) {
	owner, err := s.ownerRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrOwnerNotFound) {
			return nil, domain.ErrOwnerNotFound
		}
		return nil, err
	}
	return mapper.ToOwnerResponseDTO(owner), nil
}

func (s *OwnerService) GetByEmail(ctx context.Context, email string) (*dto.OwnerResponseDTO, error) {
	owner, err := s.ownerRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			return nil, domain.ErrEmailExists
		}
		return nil, err
	}
	return mapper.ToOwnerResponseDTO(owner), nil
}

func (s *OwnerService) GetAll(ctx context.Context)([]*dto.OwnerResponseDTO, error) {
	owners, err := s.ownerRepo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return mapper.ToOwnerListResponseDTO(owners), nil
}

func (s *OwnerService) Update(ctx context.Context, id uint, req *dto.OwnerUpdateDTO) (*dto.OwnerResponseDTO, error) {
	owner, err := s.ownerRepo.GetByID(ctx, id)

	if err != nil {
		if errors.Is(err, domain.ErrOwnerNotFound) {
			return nil, domain.ErrOwnerNotFound
		}
		return nil, err
	}

	mapper.UpdateOwnerDomain(owner, req)

	if err := s.ownerRepo.Update(ctx, owner); err != nil {
		return nil, err
	}

	return mapper.ToOwnerResponseDTO(owner), nil
}

func (s *OwnerService) Delete(ctx context.Context, id uint) error {
	if err := s.ownerRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrOwnerNotFound){
			return domain.ErrOwnerNotFound
		}
		return err
	}
	return nil
}