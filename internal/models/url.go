package models

import "time"

type URL struct {
	ID          uint   `gorm:"primaryKey"`
	OriginalURL string `gorm:"not null"`
	ShortCode   string `gorm:"uniqueIndex;not null"`
	CreatedAt   time.Time
}
