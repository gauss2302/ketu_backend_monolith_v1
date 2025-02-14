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
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type OwnerAuthService struct {
	ownerRepo    repository.OwnerRepository
	tokenManager *token.TokenManager
	cfg         *configs.JWTConfig
}

func NewOwnerAuthService(ownerRepo repository.OwnerRepository, redis *redisClient.Client, cfg *configs.JWTConfig) *OwnerAuthService {
	return &OwnerAuthService{
		ownerRepo:    ownerRepo,
		tokenManager: token.NewTokenManager(redis, cfg),
		cfg:         cfg,
	}
}

func (s *OwnerAuthService) Register(ctx context.Context, req *dto.OwnerRegisterRequestDTO) (*dto.OwnerAuthResponseDTO, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password hashing error: %w", err)
	}

	owner := &domain.Owner{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.ownerRepo.Create(ctx, owner); err != nil {
		if isPgUniqueViolation(err) {
			return nil, domain.ErrEmailExists
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	claims := token.BaseClaims{
		ID:    owner.ID,
		Email: owner.Email,
		Role:  "owner",
		Type:  "owner",
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

	return &dto.OwnerAuthResponseDTO{
		Owner:       mapper.ToOwnerResponseDTO(owner),
		AccessToken: accessToken,
		ExpiresIn:   time.Now().Add(s.cfg.AccessTTL).Unix(),
	}, nil
}

func (s *OwnerAuthService) Login(ctx context.Context, req *dto.OwnerLoginRequestDTO) (*dto.OwnerAuthResponseDTO, error) {
	log.Printf("Attempting login for email: %s", req.Email)
	
	owner, err := s.ownerRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("Error getting owner by email: %v", err)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(owner.Password), []byte(req.Password)); err != nil {
		log.Printf("Password comparison failed for email: %s", req.Email)
		return nil, domain.ErrInvalidCredentials
	}

	claims := token.BaseClaims{
		ID:    owner.ID,
		Email: owner.Email,
		Role:  "owner",
		Type:  "owner",
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(claims)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(claims)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return nil, err
	}

	err = s.tokenManager.StoreRefreshToken(ctx, claims, refreshToken)
	if err != nil {
		log.Printf("Error storing refresh token: %v", err)
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &dto.OwnerAuthResponseDTO{
		Owner:       mapper.ToOwnerResponseDTO(owner),
		AccessToken: accessToken,
		ExpiresIn:   time.Now().Add(s.cfg.AccessTTL).Unix(),
	}, nil
}

func (s *OwnerAuthService) RefreshToken(ctx context.Context, ownerID uint) (*dto.TokenRefreshResponse, error) {
	claims := token.BaseClaims{
		ID:   ownerID,
		Type: "owner",
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

func (s *OwnerAuthService) Logout(ctx context.Context, ownerID uint) error {
	claims := token.BaseClaims{
		ID:   ownerID,
		Type: "owner",
	}
	return s.tokenManager.DeleteRefreshToken(ctx, claims)
}

