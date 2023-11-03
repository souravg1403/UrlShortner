package main

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	RedisClient *redisClient
}

func RedisConnectionHandler() RedisConnection {
	connection := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	handler := RedisConnection{
		RedisClient: connection,
	}

	redis
}

func (r RedisConnection) setKey(key, value string) {
	ctx := context.Background()
	SetResult := r.RedisClient.Set(ctx, key, value)
	confirmation, err := SetResult.Result()
	if err != nil {
		return err
	} else {
		return confirmation
	}
}
