// internal/service/user.go
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

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context, req *dto.UserCreateDTO) (*dto.UserResponseDTO, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := mapper.ToUserDomain(req)
	user.Password = string(hashedPassword)

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			return nil, domain.ErrEmailExists
		}
		return nil, err
	}

	return mapper.ToUserResponseDTO(user), nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*dto.UserResponseDTO, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return mapper.ToUserResponseDTO(user), nil
}

func (s *UserService) GetAll(ctx context.Context) ([]*dto.UserResponseDTO, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.ToUserListResponseDTO(users), nil
}

func (s *UserService) Update(ctx context.Context, id uint, req *dto.UserUpdateDTO) (*dto.UserResponseDTO, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	mapper.UpdateUserDomain(user, req)

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return mapper.ToUserResponseDTO(user), nil
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}
