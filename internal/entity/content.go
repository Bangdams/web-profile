package entity

type Content struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Image       string
	Address     string
	ContactInfo string
	Category    string `gorm:"not null"`
	CreatedBy   uint   `gorm:"not null"`
	Admin       Admin  `gorm:"foreignKey:created_by;references:id"`
}
