package entity

import (
	"time"
)

type RefreshToken struct {
	AdminId      uint      `gorm:"primaryKey"`
	Token        string    `gorm:"not null"`
	StatusLogout uint      `gorm:"not null "`
	ExpiresAt    time.Time `gorm:"not null"`
	CreatedAt    time.Time
	Admin        Admin `gorm:"foreignKey:admin_id;references:id"`
}
