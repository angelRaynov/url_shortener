package cache

import (
	"context"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"time"
	"url_shortener/config"
)

type Cache struct {
	Client *redis.Client
	cfg    *config.Application
	Ctx    context.Context
}

func NewURLCache(cfg *config.Application, logger *zap.SugaredLogger) *Cache {
	c := redis.NewClient(&redis.Options{
		Addr:        cfg.RedisHost,
		DialTimeout: 5 * time.Second,
		PoolSize:    10,
	})

	if err := c.Ping(context.Background()).Err(); err != nil {
		logger.Fatalf("creating cache client:%v", err)

	}

	logger.Debug("cache started")

	return &Cache{Client: c, Ctx: context.Background()}
}
