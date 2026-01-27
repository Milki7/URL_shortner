package handlers

import (
	"net/http"
	"os"
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
		Alias   string `json:"alias"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var shortCode string
	if input.Alias != "" {
		shortCode = input.Alias
		var existing models.URL
		if err := h.DB.Where("short_code = ?", shortCode).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Alias already taken"})
			return
		}
	} else {
		for {
			shortCode = utils.GenerateRandomCode(6)
			var existing models.URL
			if err := h.DB.Where("short_code = ?", shortCode).First(&existing).Error; err != nil {
				break
			}
		}
	}

	urlEntry := models.URL{OriginalURL: input.LongURL, ShortCode: shortCode}
	h.DB.Create(&urlEntry)

	h.Redis.Set(c.Request.Context(), shortCode, input.LongURL, 24*time.Hour)

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "http://localhost:8080"
	}

	c.JSON(http.StatusOK, gin.H{"short_url": domain + "/" + shortCode})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	ctx := c.Request.Context()

	h.Redis.Incr(ctx, "clicks:"+code)

	val, err := h.Redis.Get(ctx, code).Result()
	if err == nil {
		c.Redirect(http.StatusMovedPermanently, val)
		return
	}

	var urlEntry models.URL
	if err := h.DB.Where("short_code = ?", code).First(&urlEntry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	h.Redis.Set(ctx, code, urlEntry.OriginalURL, 24*time.Hour)
	c.Redirect(http.StatusMovedPermanently, urlEntry.OriginalURL)
}

func (h *URLHandler) GetStats(c *gin.Context) {
	code := c.Param("code")
	val, err := h.Redis.Get(c.Request.Context(), "clicks:"+code).Result()
	if err != nil {
		val = "0"
	}

	c.JSON(http.StatusOK, gin.H{"short_code": code, "clicks": val})
}
