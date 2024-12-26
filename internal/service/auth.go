// internal/service/auth.go
package service

import (
	"context"
	configs "ketu_backend_monolith_v1/internal/config"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/mapper"
	repository "ketu_backend_monolith_v1/internal/repository/interfaces"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequestDTO) (*dto.AuthResponseDTO, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if isPgUniqueViolation(err) {
			return nil, domain.ErrEmailExists
		}
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.generateToken(user, s.cfg.AccessTTL, s.cfg.AccessSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user, s.cfg.RefreshTTL, s.cfg.RefreshSecret)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponseDTO{
		User:         mapper.ToUserResponseDTO(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    time.Now().Add(s.cfg.AccessTTL).Unix(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequestDTO) (*dto.AuthResponseDTO, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := s.generateToken(user, s.cfg.AccessTTL, s.cfg.AccessSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user, s.cfg.RefreshTTL, s.cfg.RefreshSecret)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponseDTO{
		User:         mapper.ToUserResponseDTO(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    time.Now().Add(s.cfg.AccessTTL).Unix(),
	}, nil
}

func (s *AuthService) generateToken(user *domain.User, ttl time.Duration, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(ttl).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func isPgUniqueViolation(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}
