package main

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	RedisClient *redis.Client
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

	return handler
}

func (r RedisConnection) setKey(key, value string) (string, error) {
	ctx := context.Background()
	SetResult := r.RedisClient.Set(ctx, key, value, 0)
	confirmation, err := SetResult.Result()
	if err != nil {
		return "failed", err
	} else {
		return confirmation, nil
	}
}

func (r RedisConnection) getValue(key string) (string, error) {
	ctx := context.Background()
	GetResult := r.RedisClient.Get(ctx, key)

	if GetResult.Val() == "" {
		return "", errors.New("key value not found in redis")
	} else {
		return GetResult.Val(), nil
	}
}

func (r RedisConnection) isExist(key string) bool {
	ctx := context.Background()
	KeyExist, _ := r.RedisClient.Exists(ctx, key).Result()

	if KeyExist == 1 {
		return true
	} else {
		return false
	}
}
