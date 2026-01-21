package models

import "time"

type URL struct {
	ID          uint   `gorm:"primarykey"`
	OriginalURL string `gorm:"notnull"`
	ShortCode   string `gorm:"uniqueIndex"`
	CreatedAt   time.Time
}
