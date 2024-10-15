package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: ", err)
	}

	return &RedisClient{rdb}
}

func (r *RedisClient) SetCache(key string, value string, ttl time.Duration) error {
	err := r.Client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetCache(key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Println("Cache miss, no value found for key:", key)
		return "", nil
	} else if err != nil {
		log.Printf("Error getting value from cache for key %s: %v", key, err)
		return "", err
	}
	log.Println("Cache hit for key:", key)
	return val, nil
}
