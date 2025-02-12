package token

import (
	"context"
	"fmt"
	"time"

	configs "ketu_backend_monolith_v1/internal/config"
	redisClient "ketu_backend_monolith_v1/internal/pkg/redis"

	"github.com/golang-jwt/jwt/v4"
)

type TokenManager struct {
	redis *redisClient.Client
	cfg   *configs.JWTConfig
}

// BaseClaims contains common claims for all token types
type BaseClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role,omitempty"`
	Type  string `json:"type"` // "user" or "owner"
}

// TokenClaims represents the structure for JWT claims
type TokenClaims struct {
	jwt.RegisteredClaims
	BaseClaims
}

func NewTokenManager(redis *redisClient.Client, cfg *configs.JWTConfig) *TokenManager {
	return &TokenManager{
		redis: redis,
		cfg:   cfg,
	}
}

func (tm *TokenManager) GenerateAccessToken(claims BaseClaims) (string, error) {
	return tm.generateToken(claims, tm.cfg.AccessTTL, tm.cfg.AccessSecret)
}

func (tm *TokenManager) GenerateRefreshToken(claims BaseClaims) (string, error) {
	return tm.generateToken(claims, tm.cfg.RefreshTTL, tm.cfg.RefreshSecret)
}

func (tm *TokenManager) generateToken(claims BaseClaims, ttl time.Duration, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		BaseClaims: claims,
	})

	return token.SignedString([]byte(secret))
}

func (tm *TokenManager) ValidateAccessToken(tokenString string) (*BaseClaims, error) {
	return tm.validateToken(tokenString, tm.cfg.AccessSecret)
}

func (tm *TokenManager) ValidateRefreshToken(tokenString string) (*BaseClaims, error) {
	return tm.validateToken(tokenString, tm.cfg.RefreshSecret)
}

func (tm *TokenManager) validateToken(tokenString string, secret string) (*BaseClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return &claims.BaseClaims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// StoreRefreshToken stores the refresh token in Redis with a prefix based on the user type
func (tm *TokenManager) StoreRefreshToken(ctx context.Context, claims BaseClaims, token string) error {
	key := fmt.Sprintf("%s:%d", claims.Type, claims.ID)
	return tm.redis.StoreRefreshToken(ctx, key, token, tm.cfg.RefreshTTL)
}

// GetStoredRefreshToken retrieves the refresh token from Redis using the appropriate prefix
func (tm *TokenManager) GetStoredRefreshToken(ctx context.Context, claims BaseClaims) (string, error) {
	key := fmt.Sprintf("%s:%d", claims.Type, claims.ID)
	return tm.redis.GetRefreshToken(ctx, key)
}

// DeleteRefreshToken removes the refresh token from Redis using the appropriate prefix
func (tm *TokenManager) DeleteRefreshToken(ctx context.Context, claims BaseClaims) error {
	key := fmt.Sprintf("%s:%d", claims.Type, claims.ID)
	return tm.redis.DeleteRefreshToken(ctx, key)
} 