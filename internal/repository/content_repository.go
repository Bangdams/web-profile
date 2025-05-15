package repository

import (
	"github.com/Bangdams/web-profile-API/internal/entity"
	"gorm.io/gorm"
)

type ContentRepository interface {
	Create(tx *gorm.DB, content *entity.Content) error
	Update(tx *gorm.DB, content *entity.Content) error
	Delete(tx *gorm.DB, content *entity.Content) error
	FindAll(tx *gorm.DB, contents *[]entity.Content) error
	FindById(tx *gorm.DB, content *entity.Content) error
	FindByIdWithAdmin(tx *gorm.DB, content *entity.Content) error
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
func (repository *ContentRepositoryImpl) FindAll(tx *gorm.DB, contents *[]entity.Content) error {
	return tx.Joins("Admin").Find(contents).Error
}

// FindById implements ContentRepository.
func (repository *ContentRepositoryImpl) FindById(tx *gorm.DB, content *entity.Content) error {
	return tx.Joins("Admin").First(content).Error
}
