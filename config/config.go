package config

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func LoadConfig() error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	return nil
}

func InitRedis(ctx context.Context) (*redis.Client, error) {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return RedisClient, nil
}
