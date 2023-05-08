package cache

import (
	"context"
	"github.com/go-redis/redis/v9"
	"log"
	"time"
	"url_shortener/config"
)

type redisCache struct {
	client *redis.Client
	cfg *config.Application
	ctx context.Context
}

type GetCacher interface {
	GetShort(long string) (string, error)
	GetLong(short string) (string, error)
	Cache(short,long string) error
}

func NewURLCache(cfg *config.Application) GetCacher {
	c := redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost,
		DialTimeout: 5 * time.Second,
		PoolSize: 10,
	})

	if err := c.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}
	return &redisCache{client: c, ctx: context.Background()}
}

func (r *redisCache) GetShort(long string) (string, error) {
	return r.client.Get(r.ctx, long).Result()
}

func (r *redisCache) setShort(short,long string) error {
	return r.client.Set(r.ctx, short, long, 24*time.Hour).Err()
}

func (r *redisCache) GetLong(short string) (string, error) {
	return r.client.Get(r.ctx, short).Result()
}

func (r *redisCache) setLong(short,long string) error {
	return r.client.Set(r.ctx, long,short , 24*time.Hour).Err()
}

func (r *redisCache) Cache(short,long string) error {
	// Create a new pipeline
	pipe := r.client.TxPipeline()

	// Execute the pipeline within a transaction
	_, err := pipe.TxPipelined(r.ctx, func(pipeliner redis.Pipeliner) error {
		pipe.Set(r.ctx, long, short, 24*time.Hour)
		pipe.Set(r.ctx, short, long, 24*time.Hour)
		return nil
	})

	return err
}