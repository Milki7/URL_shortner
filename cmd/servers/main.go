package main

import (
	"context"
	"log"

	"github.com/Milki7/URL_shortner/internal/handlers"
	"github.com/Milki7/URL_shortner/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9" // ADD THIS
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. Setup DB
	db, err := gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	db.AutoMigrate(&models.URL{})

	// 2. Setup Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Default Fedora Redis port
	})

	// Check Redis connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Could not connect to Redis. Did you start it with 'sudo systemctl start redis'?")
	}

	// 3. Setup Gin & Pass DB AND Redis to Handler
	r := gin.Default()
	handler := &handlers.URLHandler{
		DB:    db,
		Redis: rdb, // PASS REDIS HERE
	}

	r.POST("/shorten", handler.Shorten)
	r.GET("//:code", handler.Redirect)

	r.Run(":8080")
}
