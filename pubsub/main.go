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

	go subscriber(redisClient, "channel1")
	time.Sleep(1 * time.Second)
	go publisher(redisClient, "channel1")

	select {}
}

func subscriber(redisClient *redis.Client, channel string) {
	pubsub := redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()

	fmt.Println("subscribed to channel:", channel)

	for message := range pubsub.Channel() {
		fmt.Printf("received message: '%s' from channel: '%s'\n", message.Payload, channel)
	}
}

func publisher(redisClient *redis.Client, channel string) {
	messages := []string{
		"hello",
		"this is cool",
		"so is this",
	}

	for _, message := range messages {
		err := redisClient.Publish(ctx, channel, message).Err()
		if err != nil {
			log.Fatalln("could not publish message into channel:", err)
		}

		fmt.Printf("published message: '%s' in channel: '%s'\n", message, channel)
		time.Sleep(2 * time.Second)

	}
}
