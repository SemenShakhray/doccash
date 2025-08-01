package cache

import (
	"context"
	"log"
	"net"
	"strconv"

	"github.com/SemenShakhray/doccash/internal/config"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
}

func NewCache(cfg *config.Config) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: net.JoinHostPort(cfg.Cache.RedisHost, strconv.Itoa(cfg.Cache.RedisPort)),
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return &Cache{
		Client: rdb,
	}
}
