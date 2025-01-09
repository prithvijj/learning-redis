package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	counterKey = "counterkey"
)

func main() {
	pongResult, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("could not connect to redis", err)
	}
	fmt.Println("pong result:", pongResult)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", renderPage)

	router.POST("/increment", incrementCounter)

	router.Run(":8080")

}

func renderPage(c *gin.Context) {
	counter, err := getCounter()
	if err != nil {
		log.Println("error when retrieving counter value:", err)
		c.String(http.StatusInternalServerError, "could not retrieve render page")
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Counter": counter,
	})
}

func incrementCounter(c *gin.Context) {
	err := redisClient.Incr(ctx, counterKey).Err()
	if err != nil {
		c.String(http.StatusInternalServerError, "error incrementing storage")
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func getCounter() (int, error) {
	val, err := redisClient.Get(ctx, counterKey).Result()
	if err == redis.Nil {
		err = redisClient.Set(ctx, counterKey, 0, 0).Err()
		if err != nil {
			return 0, err
		}
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	counter, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return counter, nil
}
