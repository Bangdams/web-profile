package entity

type Admin struct {
	ID            uint           `gorm:"primaryKey"`
	Name          string         `gorm:"not null"`
	Username      string         `gorm:"not null;unique"`
	Password      string         `gorm:"not null"`
	Contents      []Content      `gorm:"foreignKey:created_by;references:id"`
	Announcements []Announcement `gorm:"foreignKey:published_by;references:id"`
}
