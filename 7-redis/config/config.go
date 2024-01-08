package config

import "time"

type Redis struct {
	RedisAddr    string
	MinIdleConns int
	PoolSize     int
	PoolTimeout  time.Duration
	Password     string
	DB           int
}

type Config struct {
	Redis Redis
}

func NewConfig() Config {
	return Config{
		Redis: Redis{
			RedisAddr:    "localhost:6379",
			MinIdleConns: 200,
			PoolSize:     12000,
			PoolTimeout:  240,
			Password:     "",
			DB:           0,
		},
	}
}
