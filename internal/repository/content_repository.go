package repository

import (
	"github.com/Bangdams/web-profile-API/internal/entity"
	"gorm.io/gorm"
)

type ContentRepository interface {
	Create(tx *gorm.DB, content *entity.Content) error
	Update(tx *gorm.DB, content *entity.Content) error
	Delete(tx *gorm.DB, content *entity.Content) error
	FindAll(tx *gorm.DB, order string, category string, contents *[]entity.Content) error
	FindById(tx *gorm.DB, content *entity.Content) error
}

type ContentRepositoryImpl struct {
	Repository[entity.Content]
}

func NewContentRepository() ContentRepository {
	return &ContentRepositoryImpl{}
}

// FindByIdWithAdmin implements ContentRepository.
func (repository *ContentRepositoryImpl) FindByIdWithAdmin(tx *gorm.DB, content *entity.Content) error {
	return tx.First(content).Joins("Admin").Error
}

// FindAll implements ContentRepository.
func (repository *ContentRepositoryImpl) FindAll(tx *gorm.DB, order string, category string, contents *[]entity.Content) error {
	if order == "ASC" {
		if category != "" {
			return tx.Joins("Admin").
				Order("contents.created_at ASC").
				Where("contents.category = ?", category).
				Find(contents).Error
		}
		return tx.Joins("Admin").
			Order("contents.created_at ASC").
			Find(contents).Error
	}

	if category != "" {
		return tx.Joins("Admin").
			Order("contents.created_at DESC").
			Where("contents.category = ?", category).
			Find(contents).Error
	}

	return tx.Joins("Admin").
		Order("contents.created_at DESC").
		Find(contents).Error
}

// FindById implements ContentRepository.
func (repository *ContentRepositoryImpl) FindById(tx *gorm.DB, content *entity.Content) error {
	return tx.Joins("Admin").First(content).Error
}
