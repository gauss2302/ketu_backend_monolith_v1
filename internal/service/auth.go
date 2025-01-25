// internal/service/auth.go
package service

import (
	"context"
	"fmt"
	configs "ketu_backend_monolith_v1/internal/config"
	"ketu_backend_monolith_v1/internal/domain"
	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/mapper"
	repository "ketu_backend_monolith_v1/internal/repository/interfaces"
	"log"
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
	// Add timeout context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("GetByEmail error: %v", err)
		return nil, domain.ErrInvalidCredentials
	}

	// Use constant time comparison for password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := s.generateToken(user, s.cfg.AccessTTL, s.cfg.AccessSecret)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(user, s.cfg.RefreshTTL, s.cfg.RefreshSecret)
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	expiresIn := time.Now().Add(s.cfg.AccessTTL).Unix()

	return &dto.AuthResponseDTO{
		User:         mapper.ToUserResponseDTO(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
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

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponseDTO, error) {
    // Parse and validate refresh token
    token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(s.cfg.RefreshSecret), nil
    })

    if err != nil {
        return nil, domain.ErrInvalidCredentials
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, domain.ErrInvalidCredentials
    }

    // Get user from database
    userID := uint(claims["user_id"].(float64))
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, domain.ErrInvalidCredentials
    }

    // Generate new tokens
    accessToken, err := s.generateToken(user, s.cfg.AccessTTL, s.cfg.AccessSecret)
    if err != nil {
        return nil, fmt.Errorf("failed to generate access token: %w", err)
    }

    newRefreshToken, err := s.generateToken(user, s.cfg.RefreshTTL, s.cfg.RefreshSecret)
    if err != nil {
        return nil, fmt.Errorf("failed to generate refresh token: %w", err)
    }

    return &dto.AuthResponseDTO{
        User:         mapper.ToUserResponseDTO(user),
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken,
        ExpiresIn:    time.Now().Add(s.cfg.AccessTTL).Unix(),
    }, nil
}

func isPgUniqueViolation(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == "23505"
	}
	return false
}
