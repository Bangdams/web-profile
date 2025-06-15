package entity

import "time"

type Announcement struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Content     string `gorm:"not null"`
	Image       string
	PublishedBy uint `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Admin       Admin `gorm:"foreignKey:published_by;references:id"`
}
