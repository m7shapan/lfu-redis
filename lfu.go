package lfu

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const LFUHashName = "LFUCacheHash"
const LFUSortedSetName = "LFUCacheSortedSet"

type LFU struct {
	redisClient *redis.Client
	capacity    int64
}

func New(capacity int64, redisClient *redis.Client) *LFU {
	return &LFU{
		redisClient: redisClient,
		capacity:    capacity,
	}
}

func (c *LFU) Put(key string, value interface{}) (err error) {
	capacity, _ := c.redisClient.ZCard(context.Background(), LFUSortedSetName).Result()
	if capacity >= c.capacity {
		deleteCount := capacity - c.capacity + 1
		items := c.redisClient.ZPopMin(context.Background(), LFUSortedSetName, deleteCount).Val()

		var keys []string
		for _, item := range items {
			keys = append(keys, item.Member.(string))
		}

		c.redisClient.HDel(context.Background(), LFUHashName, keys...)
	}

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
