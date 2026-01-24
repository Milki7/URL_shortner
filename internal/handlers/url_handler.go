package handlers

import (
	"net/http"

	"time"

	"github.com/Milki7/URL_shortner/internal/models"
	"github.com/Milki7/URL_shortner/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type URLHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var input struct {
		LongURL string `json:"long_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var shortCode string
	// Collision Handling: Loop until we find a unique code
	for {
		shortCode = utils.GenerateRandomCode(6) // 6 chars = ~56 billion combinations

		var existing models.URL
		// Check if this code is already in the DB
		result := h.DB.Where("short_code = ?", shortCode).First(&existing)
		if result.Error != nil {
			// If error is "Record Not Found", the code is unique!
			break
		}
	}

	urlEntry := models.URL{
		OriginalURL: input.LongURL,
		ShortCode:   shortCode,
	}

	if err := h.DB.Create(&urlEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save URL"})
		return
	}

	// Warm the Redis cache immediately for speed
	h.Redis.Set(c.Request.Context(), shortCode, input.LongURL, 24*time.Hour)

	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8080/" + shortCode})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	ctx := c.Request.Context()

	// 1. Try to get the URL from Redis (Fast Path)
	val, err := h.Redis.Get(ctx, code).Result()
	if err == nil {
		c.Redirect(http.StatusMovedPermanently, val)
		return
	}

	// 2. Cache Miss: Look in PostgreSQL/SQLite (Slow Path)
	var urlEntry models.URL
	if err := h.DB.Where("short_code = ?", code).First(&urlEntry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// 3. Store in Redis for next time (with a 24-hour expiration)
	h.Redis.Set(ctx, code, urlEntry.OriginalURL, 24*time.Hour)

	c.Redirect(http.StatusMovedPermanently, urlEntry.OriginalURL)
}
