// internal/service/auth.go
package service

import (
	"context"
	"fmt"
	configs "ketu_backend_monolith_v1/internal/config"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/mapper"
	redisClient "ketu_backend_monolith_v1/internal/pkg/redis"
	"ketu_backend_monolith_v1/internal/pkg/utils/token"
	repository "ketu_backend_monolith_v1/internal/repository/interfaces"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     repository.UserRepository
	tokenManager *token.TokenManager
	cfg         *configs.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, redis *redisClient.Client, cfg *configs.JWTConfig) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenManager: token.NewTokenManager(redis, cfg),
		cfg:         cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequestDTO) (*dto.AuthResponseDTO, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password hashing error: %w", err)
	}

	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Name:      req.Name,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if isPgUniqueViolation(err) {
			return nil, domain.ErrEmailExists
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	claims := token.BaseClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  "user",
		Type:  "user",
	}

	// Generate tokens
	accessToken, err := s.tokenManager.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(claims)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	err = s.tokenManager.StoreRefreshToken(ctx, claims, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &dto.AuthResponseDTO{
		User:        mapper.ToUserResponseDTO(user),
		AccessToken: accessToken,
		ExpiresIn:   time.Now().Add(s.cfg.AccessTTL).Unix(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequestDTO) (*dto.AuthResponseDTO, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	claims := token.BaseClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  "user",
		Type:  "user",
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(claims)
	if err != nil {
		return nil, err
	}

	err = s.tokenManager.StoreRefreshToken(ctx, claims, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &dto.AuthResponseDTO{
		User:        mapper.ToUserResponseDTO(user),
		AccessToken: accessToken,
		ExpiresIn:   time.Now().Add(s.cfg.AccessTTL).Unix(),
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, userID uint) (*dto.TokenRefreshResponse, error) {
	claims := token.BaseClaims{
		ID:   userID,
		Type: "user",
	}
	
	storedToken, err := s.tokenManager.GetStoredRefreshToken(ctx, claims)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	tokenClaims, err := s.tokenManager.ValidateRefreshToken(storedToken)
	if err != nil {
		_ = s.tokenManager.DeleteRefreshToken(ctx, claims)
		return nil, err
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(*tokenClaims)
	if err != nil {
		return nil, err
	}

	return &dto.TokenRefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   int64(s.cfg.AccessTTL.Seconds()),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uint) error {
	claims := token.BaseClaims{
		ID:   userID,
		Type: "user",
	}
	return s.tokenManager.DeleteRefreshToken(ctx, claims)
}

func isPgUniqueViolation(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == "23505"
	}
	return false
}
