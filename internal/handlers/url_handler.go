package handlers

import (
	"net/http"

	"time"

	"github.com/Milki7/URL_shortner/internal/models"
	"github.com/Milki7/URL_shortner/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type URLHandler struct {
	DB *gorm.DB
}

// Shorten handles POST /shorten
func (h *URLHandler) Shorten(c *gin.Context) {
	var input struct {
		LongURL string `json:"long_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urlEntry := models.URL{OriginalURL: input.LongURL}
	h.DB.Create(&urlEntry) // Create entry to get the ID

	// Generate code from ID and update record
	shortCode := utils.Encode(urlEntry.ID)
	h.DB.Model(&urlEntry).Update("ShortCode", shortCode)

	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8080/" + shortCode})
}

// Redirect handles GET /:code
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
