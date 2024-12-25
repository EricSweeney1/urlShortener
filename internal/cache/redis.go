package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Dashboard/urlShortener/config"
	"github.com/Dashboard/urlShortener/internal/repo"
	"github.com/go-redis/redis/v8"
	"time"
)

//SetURL(ctx context.Context, url repo.Url) error

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(config config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RedisCache{client: client}, nil
}

func (c *RedisCache) SetURL(ctx context.Context, url repo.Url) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}
	if err := c.client.Set(ctx, url.ShortCode, data, time.Until(url.ExpiredTime)).Err(); err != nil {
		return err
	}
	return nil
}
func (c *RedisCache) GetURL(ctx context.Context, shortCode string) (*repo.Url, error) {
	data, err := c.client.Get(ctx, shortCode).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var url repo.Url
	if err := json.Unmarshal(data, &url); err != nil {
		return nil, err
	}
	return &url, nil

}
func (c *RedisCache) Close() error {
	return c.client.Close()
}
