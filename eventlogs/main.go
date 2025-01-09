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

	eventLogKey := "eventLogs"

	logEvent(redisClient, eventLogKey, "Alpha user joined")
	logEvent(redisClient, eventLogKey, "Beta user joined")
	logEvent(redisClient, eventLogKey, "Alpha user changed something")
	logEvent(redisClient, eventLogKey, "Alpha user exited")
	logEvent(redisClient, eventLogKey, "Beta user exited")

	recentEvents, err := getRecentEvents(redisClient, eventLogKey, 3)
	if err != nil {
		log.Fatalln("could not get recent events:", err)
	}

	for i, event := range recentEvents {
		fmt.Printf("idx: %v, log: %s\n", i+1, event)
	}

}

func logEvent(redisClient *redis.Client, eventLogKey string, logMessage string) {
	timeStampEvent := fmt.Sprintf("%s: %s", time.Now().Format(time.RFC3339), logMessage)
	err := redisClient.RPush(ctx, eventLogKey, timeStampEvent).Err()
	if err != nil {
		log.Fatalln("could not push log event:", err)
	}

}

func getRecentEvents(redisClient *redis.Client, eventLogKey string, count int64) ([]string, error) {
	return redisClient.LRange(ctx, eventLogKey, -count, -1).Result()
}
