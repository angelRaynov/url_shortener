package repository

import (
	"github.com/go-redis/redis/v9"
	"time"
	"url_shortener/infrastructure/cache"
)

type cacheRepo struct {
	cache *cache.Cache
}

func NewCacheRepo(c *cache.Cache) *cacheRepo {
	return &cacheRepo{cache: c}
}

func (cr *cacheRepo) GetShort(long string) (string, error) {
	return cr.cache.Client.Get(cr.cache.Ctx, long).Result()
}

func (cr *cacheRepo) setShort(short, long string) error {
	return cr.cache.Client.Set(cr.cache.Ctx, short, long, 24*time.Hour).Err()
}

func (cr *cacheRepo) GetLong(short string) (string, error) {
	return cr.cache.Client.Get(cr.cache.Ctx, short).Result()
}

func (cr *cacheRepo) setLong(short, long string) error {
	return cr.cache.Client.Set(cr.cache.Ctx, long, short, 24*time.Hour).Err()
}

func (cr *cacheRepo) Cache(short, long string) error {
	// Create a new pipeline
	pipe := cr.cache.Client.TxPipeline()

	// Execute the pipeline within a transaction
	_, err := pipe.TxPipelined(cr.cache.Ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(cr.cache.Ctx, long, short, 24*time.Hour)
		pipe.Set(cr.cache.Ctx, short, long, 24*time.Hour)
		return nil
	})

	return err
}
