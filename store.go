package main

import "github.com/go-redis/redis"

// Global redis client
var rdb = redis.NewClient(&redis.Options{
	Addr:     getEnv("REDISHOST", "localhost") + ":" + getEnv("REDISPORT", "6379"),
	Password: getEnv("REDISPASSWORD", ""),
	DB:       0, // use default DB
})
