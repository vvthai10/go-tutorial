package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/vvthai10/redis/config"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	redisHost := cfg.Redis.RedisAddr
	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
		Password:     cfg.Redis.Password, // no password set
		DB:           cfg.Redis.DB,       // use default DB
	})

	return &RedisClient{
		client: client,
	}
}

func (r *RedisClient) Get(key string) ([]byte, error) {
	userBytes, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "RedisClient.Get.client.Get")
	}
	return userBytes, nil
}
func (r *RedisClient) Set(key string, seconds int, stored interface{}) error {
	storedBytes, err := json.Marshal(stored)
	if err != nil {
		return errors.Wrap(err, "RedisClient.Set.json.Marshal")
	}
	if err := r.client.Set(key, storedBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "RedisClient.Set.client.Set")
	}
	return nil
}

func (r *RedisClient) SetExpires(key string, seconds int) error {
	err := r.client.Expire(key, time.Second*time.Duration(seconds)).Err()
	if err != nil {
		return errors.Wrap(err, "RedisClient.SetExpire.client.Expire")
	}
	return nil
}

func (r *RedisClient) Remove(key string) error {
	if err := r.client.Del(key).Err(); err != nil {
		return errors.Wrap(err, "RedisClient.Remove.client.Del")
	}
	return nil
}
