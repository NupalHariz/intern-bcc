package redis

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func ConnectToRedis() {
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal("Failed to Connect to Redis")
	}

	redisClient := redis.NewClient(opts)
	RDB = redisClient
}

type IRedis interface {
	SetRedis(ctx context.Context, key string, data string, ttl time.Duration) error
	GetRedis(ctx context.Context, key string) (string, error)
}

type Redis struct {
	r *redis.Client
}

func RedisInit(r *redis.Client) IRedis {
	return &Redis{r}
}

func (r *Redis) SetRedis(ctx context.Context, key string, data string, ttl time.Duration) error {
	err := r.r.SetEx(ctx, key, data, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetRedis(ctx context.Context, key string) (string, error) {
	stringData, err := r.r.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return stringData, nil
}
