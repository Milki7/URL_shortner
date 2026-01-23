package models

import "time"

type URL struct {
	ID          uint   `gorm:"primaryKey"`
	OriginalURL string `gorm:"notNull"`
	ShortCode   string `gorm:"uniqueIndex"`
	CreatedAt   time.Time
}
