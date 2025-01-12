package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	ctx  = context.Background()
	port = ":8080"
)

// Initialize Redis client
func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Default DB
	})
}

// Convert JPEG image file to Base64
func imageToBase64(imgReader io.Reader) (string, error) {
	// Decode the JPEG image
	img, _, err := image.Decode(imgReader)
	if err != nil {
		return "", fmt.Errorf("error decoding image: %v", err)
	}

	// Convert image to buffer (JPEG format)
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return "", fmt.Errorf("error encoding image to JPEG: %v", err)
	}

	// Base64 encode
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// Store Base64-encoded image in Redis
func storeImageInRedis(key, base64Image string) error {
	// Store the Base64 image with a TTL of 1 hour
	return rdb.Set(ctx, key, base64Image, time.Hour).Err()
}

// Get image from Redis and decode from Base64
func getImageFromRedis(key string) (image.Image, error) {
	base64Image, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// Decode Base64 image
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return nil, fmt.Errorf("error decoding Base64 string: %v", err)
	}

	// Decode image from the byte slice
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("error decoding image from byte slice: %v", err)
	}

	return img, nil
}

// API endpoint to upload image and store it in Redis
func uploadImage(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error"})
		return
	}
	defer file.Close()

	// Convert image to Base64
	base64Image, err := imageToBase64(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store Base64 image in Redis
	key := "image_key"
	err = storeImageInRedis(key, base64Image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing image in Redis"})
		return
	}

	// Redirect to the homepage with the success message
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"message": "JPEG image uploaded and stored in Redis successfully!",
	})
}

// API endpoint to retrieve image from Redis and display it
func getImage(c *gin.Context) {
	key := "image_key" // Using the same key as used during upload

	// Get image from Redis
	img, err := getImageFromRedis(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found in Redis"})
		return
	}

	// Serve the image with the correct MIME type for JPEG
	c.Header("Content-Type", "image/jpeg")
	err = jpeg.Encode(c.Writer, img, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encoding image"})
	}
}

func main() {
	// Initialize Redis client
	initRedis()

	// Set up Gin routes
	r := gin.Default()

	// Serve HTML form at the root route
	r.LoadHTMLFiles("templates/upload.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	// Upload image
	r.POST("/upload", uploadImage)

	// Get image from Redis
	r.GET("/image", getImage)

	// Start server
	r.Run(port)
}
