package cache

import (
	"log"
	"os"

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
