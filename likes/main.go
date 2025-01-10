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
	likeSetKey = "post:likes"
)

func main() {
	pongResult, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("could not connect to redis", err)
	}
	fmt.Println("pong result:", pongResult)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", showLikesPage)

	router.POST("/like", addUserToLikes)

	router.Run(":8080")

}

func showLikesPage(c *gin.Context) {
	users, err := redisClient.SMembers(ctx, likeSetKey).Result()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "Failed to fetch likes",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"users": users,
	})
}

func addUserToLikes(c *gin.Context) {
	user := c.PostForm("user")

	if user == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "User name cannot be empty",
		})
		return
	}

	err := redisClient.SAdd(ctx, likeSetKey, user).Err()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "Failed to add users to likes",
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}
