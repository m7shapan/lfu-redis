package lfu

import "github.com/go-redis/redis/v8"

type LFU struct {
	redisClient *redis.Client
}

func New(redisClient *redis.Client) *LFU {
	return &LFU{
		redisClient: redisClient,
	}
}
