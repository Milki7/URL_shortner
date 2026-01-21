package main

import (
	"github.com/Milki7/URL_shortner/internal/handlers"
	"github.com/Milki7/URL_shortner/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. Setup DB (using SQLite for local dev on Fedora)
	db, err := gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.URL{})

	// 2. Setup Gin
	r := gin.Default()
	handler := &handlers.URLHandler{DB: db}

	// 3. Routes
	r.POST("/shorten", handler.Shorten)
	r.GET("/:code", handler.Redirect)

	r.Run(":8080")
}
