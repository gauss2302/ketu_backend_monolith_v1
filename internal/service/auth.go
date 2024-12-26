package service

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	configs "ketu_backend_monolith_v1/internal/config"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	repository "ketu_backend_monolith_v1/internal/repository/interfaces"

	"time"
)

type AuthService struct {
	userRepo repository.UserRepository
	cfg      *configs.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, cfg *configs.JWTConfig) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, input dto.RegisterRequest) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		// Check if it's a unique constraint violation
		if isPgUniqueViolation(err) {
			return nil, domain.ErrEmailExists
		}
		return nil, err
	}

	return user, nil
}

func isPgUniqueViolation(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}

func (s *AuthService) Login(ctx context.Context, input dto.LoginRequest) (*domain.User, string, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, "", "", domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, "", "", domain.ErrInvalidCredentials
	}

	accessToken, err := s.generateToken(user, s.cfg.AccessTTL, s.cfg.AccessSecret)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.generateToken(user, s.cfg.RefreshTTL, s.cfg.RefreshSecret)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) generateToken(user *domain.User, ttl time.Duration, secret string) (string, error) {
	// Change user*id to user_id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(ttl).Unix(),
	})

	return token.SignedString([]byte(secret))
}
