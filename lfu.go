package lfu

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const LFUHashName = "LFUCacheHash"
const LFUSortedSetName = "LFUCacheSortedSet"

type LFUCache struct {
	redisClient *redis.Client
	capacity    int64
}

func New(capacity int64, redisClient *redis.Client) *LFUCache {
	return &LFUCache{
		redisClient: redisClient,
		capacity:    capacity,
	}
}

// Put caches item
func (c *LFUCache) Put(key string, value interface{}) (err error) {
	c.capacityCheck()

	err = c.redisClient.HSet(context.Background(), LFUHashName, key, value).Err()

	c.frequentlyIncr(key)
	return
}

// Get return item from cache
func (c *LFUCache) Get(key string) (value string, err error) {
	value, err = c.redisClient.HGet(context.Background(), LFUHashName, key).Result()

	if err != nil {
		return
	}

	c.frequentlyIncr(key)
	return
}

// DelItem remove item from cache
func (c *LFUCache) DelItem(key string) (err error) {
	err = c.redisClient.HDel(context.Background(), LFUHashName, key).Err()
	if err != nil {
		return
	}

	err = c.redisClient.ZRem(context.Background(), LFUSortedSetName, key).Err()
	return
}

// Flush remove all cache
func (c *LFUCache) Flush() (err error) {
	err = c.redisClient.Del(context.Background(), LFUHashName).Err()
	if err != nil {
		return
	}

	err = c.redisClient.Del(context.Background(), LFUSortedSetName).Err()
	return
}

func (c *LFUCache) capacityCheck() {
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
}

func (c *LFUCache) frequentlyIncr(key string) {
	c.redisClient.ZIncr(context.Background(), LFUSortedSetName, &redis.Z{
		Member: key,
		Score:  1,
	})
}
