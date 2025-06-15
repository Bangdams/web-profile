package repository

import (
	"github.com/Bangdams/web-profile-API/internal/entity"
	"gorm.io/gorm"
)

type AnnouncementRepository interface {
	Create(tx *gorm.DB, announcement *entity.Announcement) error
	Update(tx *gorm.DB, announcement *entity.Announcement) error
	Delete(tx *gorm.DB, announcement *entity.Announcement) error
	FindAll(tx *gorm.DB, order string, announcements *[]entity.Announcement) error
	FindById(tx *gorm.DB, announcement *entity.Announcement) error
	GetFirst(tx *gorm.DB, announcement *entity.Announcement) error
}

type AnnouncementRepositoryImpl struct {
	Repository[entity.Announcement]
}

func NewAnnouncementRepository() AnnouncementRepository {
	return &AnnouncementRepositoryImpl{}
}

// FindAll implements AnnouncementRepository.
func (repository *AnnouncementRepositoryImpl) FindAll(tx *gorm.DB, order string, announcements *[]entity.Announcement) error {
	if order == "ASC" {
		return tx.Joins("Admin").
			Order("announcements.created_at ASC").
			Find(announcements).Limit(4).Error
	}
	return tx.Joins("Admin").
		Order("announcements.created_at DESC").
		Find(announcements).Limit(4).Error
}

// FindById implements AnnouncementRepository.
func (repository *AnnouncementRepositoryImpl) FindById(tx *gorm.DB, announcement *entity.Announcement) error {
	return tx.Joins("Admin").First(announcement).Error
}

// GetFirst implements AnnouncementRepository.
func (repository *AnnouncementRepositoryImpl) GetFirst(tx *gorm.DB, announcement *entity.Announcement) error {
	return tx.Joins("Admin").Order("announcements.created_at DESC").First(announcement).Error
}
