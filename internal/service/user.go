package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/handler/dto"
	"ketu_backend_monolith_v1/internal/repository"
	"time"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, input dto.CreateUserInput) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, id uint) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, id uint, input dto.UpdateUserInput) error {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	user.UpdatedAt = time.Now()
	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id uint) error {
	return uc.userRepo.Delete(ctx, id)
}

func (uc *UserUseCase) GetAll(ctx context.Context) ([]domain.User, error) {
	return uc.userRepo.GetAll(ctx)
}
