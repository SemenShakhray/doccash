package cache

import (
	"context"
	"time"
)

func (c *Cache) SetToken(ctx context.Context, token string, ttl time.Duration) error {
	return c.Client.Set(ctx, token, "valid", ttl).Err()
}

func (c *Cache) DeleteToken(ctx context.Context, token string) error {
	return c.Client.Del(ctx, token).Err()
}

func (c *Cache) IsTokenValid(ctx context.Context, token string) (bool, error) {
	result, err := c.Client.Exists(ctx, token).Result()
	return result == 1, err
}
