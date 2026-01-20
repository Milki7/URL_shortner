package models

import "time"

type URL struct {
	ID           uint   `gorm:"primarykey"`
	OringinalURL string `gorm:"notnull"`
	ShortCode    string `gorm:"uniqueIndex"`
	CreatedAt    time.Time
}
