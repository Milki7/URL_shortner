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
		Alias   string `json:"alias"` // Optional field
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var shortCode string

	// CASE 1: User provided a custom alias
	if input.Alias != "" {
		shortCode = input.Alias
		var existing models.URL
		// Check if the custom alias is already taken
		if err := h.DB.Where("short_code = ?", shortCode).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Alias already in use"})
			return
		}
	} else {
		// CASE 2: No alias provided, generate a random one
		for {
			shortCode = utils.GenerateRandomCode(6)
			var existing models.URL
			if err := h.DB.Where("short_code = ?", shortCode).First(&existing).Error; err != nil {
				break // Unique code found
			}
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

	// Cache it in Redis immediately
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
