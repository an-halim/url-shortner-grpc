package entity

import "time"

type Url struct {
	ID        int    `gorm:"primaryKey"`
	ShortUrl  string `gorm:"unique"`
	Original  string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
