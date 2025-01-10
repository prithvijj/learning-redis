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
	hashKey = "user:data"
)

func main() {

	pongResult, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("could not connect to redis", err)
	}
	fmt.Println("pong result:", pongResult)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", showHashPage)
	router.POST("/add", addKeyValueToHash)

	router.Run(":8080")
}

func showHashPage(c *gin.Context) {
	hashData, err := redisClient.HGetAll(ctx, hashKey).Result()
	if err != nil {
		log.Println("error when calling HGetAll", err)
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "could not get all hash keys",
		})
		return
	}

	fieldCount, err := redisClient.HLen(ctx, hashKey).Result()
	if err != nil {
		log.Println("error when calling HLen", err)
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "could not get length of hash keys",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"hashData":   hashData,
		"fieldCount": fieldCount,
	})
}

func addKeyValueToHash(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")

	if key == "" || value == "" {
		c.HTML(http.StatusBadGateway, "index.html", gin.H{
			"error": "the form should not be empty",
		})
		return
	}

	err := redisClient.HSet(ctx, hashKey, key, value).Err()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "could not set keys in hashes",
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}
