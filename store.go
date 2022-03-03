package main

import "github.com/go-redis/redis"

// Global redis client
var rdb = redis.NewClient(&redis.Options{
	Addr: getEnv("REDIS_URL", "localhost:6379"),
	DB:   0, // use default DB
})
