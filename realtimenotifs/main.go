package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	channel = "notifications"
)

func main() {
	pongResult, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("could not connect to redis", err)
	}
	fmt.Println("pong result:", pongResult)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", showHomePage)

	router.GET("/stream", streamHandler)

	router.POST("/notify", notifyHandler)

	router.Run(":8080")
}

func showHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func streamHandler(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	pubsub := redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(c.Writer, "data: %s\n\n", msg.Payload)
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}

	// for msg := range ch {
	// 	fmt.Println("Received message:", msg.Payload) // Add this for debugging
	// 	fmt.Fprintf(c.Writer, "data: %s\n", msg.Payload)
	// 	c.Writer.Flush()
	// }
}

func notifyHandler(c *gin.Context) {
	message := c.PostForm("message")

	if message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "message was empty",
		})
		return
	}
	err := redisClient.Publish(ctx, channel, message).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to send notifications",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "notifications sent",
	})
}
