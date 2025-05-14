package repository

import (
	"github.com/Bangdams/web-profile-API/internal/entity"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(tx *gorm.DB, refreshToken *entity.RefreshToken) error
	Update(tx *gorm.DB, refreshToken *entity.RefreshToken) error
	FindById(tx *gorm.DB, adminId uint) error
	CheckStatusLogout(tx *gorm.DB, adminId uint) error
}

type RefreshTokenRepositoryImpl struct {
	Repository[entity.RefreshToken]
}

func NewRefreshTokenRepository() RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{}
}

// CheckStatusLogout implements UserRepository.
func (repository *RefreshTokenRepositoryImpl) CheckStatusLogout(tx *gorm.DB, adminId uint) error {
	return tx.First(&entity.RefreshToken{}, "admin_id = ? AND status_logout = ?", adminId, 0).Error
}

// FindById implements UserRepository.
func (repository *RefreshTokenRepositoryImpl) FindById(tx *gorm.DB, adminId uint) error {
	return tx.First(&entity.RefreshToken{}, "admin_id=?", adminId).Error
}
