package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	// Create
	err = redisClient.Set(ctx, "key1", "value1", 10*time.Second).Err()
	if err != nil {
		log.Fatalln("could not set key:", err)
	}

	// Read
	value, err := redisClient.Get(ctx, "key1").Result()
	if err != nil {
		log.Fatalln("could not get key:", err)
	}
	fmt.Println("value retrieved was:", value)

	// Update
	err = redisClient.Set(ctx, "key1", "value2", 10*time.Second).Err()
	if err != nil {
		log.Fatalln("could not set key:", err)
	}

	value, err = redisClient.Get(ctx, "key1").Result()
	if err != nil {
		log.Fatalln("could not get key:", err)
	}
	fmt.Println("updated value retrieved was:", value)

	// Delete
	_, err = redisClient.Del(ctx, "key1").Result()
	if err != nil {
		log.Fatalln("could not delete key:", err)
	}
	fmt.Println("successfully deleted the given key")

}
