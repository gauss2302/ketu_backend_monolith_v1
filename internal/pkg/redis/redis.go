package redis

import (
	"context"
	"fmt"
	"ketu_backend_monolith_v1/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func NewRedisClient(cfg *config.RedisConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.URL,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{client}, nil
}

func (c *Client) StoreRefreshToken(ctx context.Context, key string, refreshToken string, expiration time.Duration) error {
	return c.Set(ctx, key, refreshToken, expiration).Err()
}

func (c *Client) GetRefreshToken(ctx context.Context, key string) (string, error) {
	token, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("refresh token not found")
	}
	return token, err
}

func (c *Client) DeleteRefreshToken(ctx context.Context, key string) error {
	return c.Del(ctx, key).Err()
}

func (c *Client) Close() error {
	return c.Client.Close()
} 