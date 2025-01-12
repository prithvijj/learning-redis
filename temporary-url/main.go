package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

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
)

func main() {
	pongResult, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("could not connect to redis", err)
	}
	fmt.Println("pong result:", pongResult)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")

	router.GET("/generate", showGenerateForm)
	router.POST("/generate", generateTemporaryURL)
	router.GET("/access/:key", accessKeyTemporaryURL)

	router.Run(":8080")
}

func showGenerateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "generate.html", nil)
}

func generateTemporaryURL(c *gin.Context) {
	content := c.PostForm("content")
	exp := c.PostForm("expiration")

	expDuration, err := time.ParseDuration(exp + "s")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid expiration time",
		})
		return
	}

	key := time.Now().UnixNano()

	err = redisClient.Set(ctx, strconv.Itoa(int(key)), content, expDuration).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to set content and exp",
		})
		return
	}

	fmt.Println(c.Request.Host)
	tempURL := "access/" + strconv.Itoa(int(key))
	safeURL := template.URLQueryEscaper(tempURL)

	fmt.Println(safeURL)
	c.HTML(http.StatusOK, "generate.html", gin.H{
		"message": "temporary url created",
		"urlz":    safeURL,
	})
}

func accessKeyTemporaryURL(c *gin.Context) {
	key := c.Param("key")

	content, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		c.HTML(http.StatusNotFound, "access.html", gin.H{
			"error": "content does not exist or has expired",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not retrieve content",
		})
		return
	}

	c.HTML(http.StatusOK, "access.html", gin.H{
		"content": content,
	})

}
