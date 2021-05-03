package redis

import (
	"github.com/go-redis/redis"
)

var (
	RedisClient *redis.Client
)

func init() {
	RedisClient = getRedisClient()
}

func getRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{})
}
