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

	leaderboardKey := "leaderboard1"

	for value := range 20 {
		addPlayer(redisClient, leaderboardKey, fmt.Sprintf("Player%d", value), float64(value*10))
	}

	topPlayers, err := getTopPlayers(redisClient, leaderboardKey, 10)
	if err != nil {
		log.Fatalln("could not retrieve top players:", err)
	}

	for _, topPlayer := range topPlayers {

		fmt.Printf("%v, %s\n", topPlayer.Score, topPlayer.Member)
	}
}

func addPlayer(redisClient *redis.Client, leaderboardKey, playerName string, playerScore float64) {
	err := redisClient.ZAdd(ctx, leaderboardKey, redis.Z{
		Score:  playerScore,
		Member: playerName,
	}).Err()
	if err != nil {
		log.Fatalln("could not add player:", err)
	}
}

func getTopPlayers(redisClient *redis.Client, leaderboardKey string, topN int64) ([]redis.Z, error) {
	return redisClient.ZRevRangeWithScores(ctx, leaderboardKey, 0, topN-1).Result()
}
