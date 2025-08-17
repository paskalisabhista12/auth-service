package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis(addr string, password string) error {
	db := 0 // use default DB

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		return err
	}

	log.Println("✅ Connected to Redis")
	return nil
}

func Set(key string, value string, expiration time.Duration) error {
	return Rdb.Set(Ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return Rdb.Get(Ctx, key).Result()
}

func Del(key string) error {
	return Rdb.Del(Ctx, key).Err()
}
