package main

import (
	"fmt"

	"github.com/vvthai10/redis/config"
	"github.com/vvthai10/redis/redis"
)

func main() {
	cfg := config.NewConfig()
	redisClient := redis.NewRedisClient(&cfg)
	err := redisClient.Set("demo", 5, "Storage demo")
	if err != nil {
		fmt.Println(err)
	}

	stored, err := redisClient.Get("demo")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(stored)

}
