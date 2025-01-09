package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pongResult, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("could not connect to redis", err)
	}
	fmt.Println("pong result:", pongResult)
}
