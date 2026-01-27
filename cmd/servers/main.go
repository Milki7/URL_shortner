package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Milki7/URL_shortner/internal/handlers"
	"github.com/Milki7/URL_shortner/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// 1. Setup DB
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "urls.db"
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	db.AutoMigrate(&models.URL{})

	// 2. Setup Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Redis connection failed. Run: sudo systemctl start redis")
	}

	// 3. Setup Gin
	r := gin.Default()

	// Load HTML frontend
	r.LoadHTMLFiles("web/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	handler := &handlers.URLHandler{
		DB:    db,
		Redis: rdb,
	}

	// 4. Routes
	r.POST("/shorten", handler.Shorten)
	r.GET("/:code", handler.Redirect)
	r.GET("/stats/:code", handler.GetStats)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
