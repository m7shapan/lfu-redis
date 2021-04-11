# LFU Redis cache library for Golang
LFU Redis implements LFU Cache algorithm using Redis as data storage

LFU Redis Package gives you control over Cache Capacity in case you're using multipurpose Redis instance and avoid using eviction policy

## Installation 
```bash
go get -u github.com/m7shapan/lfu-redis
```

## Quickstart
```go
package lfu_test

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/m7shapan/lfu-redis"
)

func ExampleUsage() {
	redisClient := redis.NewClient(&redis.Options{})

	cache := lfu.New(10000, redisClient)

	err := cache.Put("key", "value")
	if err != nil {
		panic(err)
	}

	value, err := cache.Get("key")
	if err != nil {
		panic(err)
	}

	fmt.Println(value)
	// Output: value
}
```