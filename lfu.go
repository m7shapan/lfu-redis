package lfu

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const LFUHashName = "LFUCacheHash"
const LFUSortedSetName = "LFUCacheSortedSet"

type LFU struct {
	redisClient *redis.Client
}

func New(redisClient *redis.Client) *LFU {
	return &LFU{
		redisClient: redisClient,
	}
}

func (c *LFU) Put(key string, value interface{}) (err error) {
	err = c.redisClient.HSet(context.Background(), LFUHashName, key, value).Err()

	c.redisClient.ZAdd(context.Background(), LFUSortedSetName, &redis.Z{
		Member: key,
		Score:  1,
	})

	return
}

func (c *LFU) Get(key string) (value string, err error) {
	value, err = c.redisClient.HGet(context.Background(), LFUHashName, key).Result()

	if err != nil {
		return
	}

	c.redisClient.ZIncr(context.Background(), LFUSortedSetName, &redis.Z{
		Member: key,
		Score:  1,
	})
	return
}
