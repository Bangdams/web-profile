package repository

import (
	"github.com/Bangdams/web-profile-API/internal/entity"
	"gorm.io/gorm"
)

type AdminRepository interface {
	Create(tx *gorm.DB, admin *entity.Admin) error
	Update(tx *gorm.DB, admin *entity.Admin) error
	Delete(tx *gorm.DB, admin *entity.Admin) error
	FindAll(tx *gorm.DB, adminId uint, admins *[]entity.Admin) error
	FindById(tx *gorm.DB, admin *entity.Admin) error
	FindByEmail(tx *gorm.DB, admin *entity.Admin) error
	Login(tx *gorm.DB, admin *entity.Admin, keyword string) error
}

type AdminRepositoryImpl struct {
	Repository[entity.Admin]
}

func NewAdminRepository() AdminRepository {
	return &AdminRepositoryImpl{}
}

// Login implements AdminRepository.
func (repository *AdminRepositoryImpl) Login(tx *gorm.DB, admin *entity.Admin, keyword string) error {
	return tx.Where("email = ?", keyword).First(admin).Error
}

// FindByEmail implements AdminRepository.
func (repository *AdminRepositoryImpl) FindByEmail(tx *gorm.DB, admin *entity.Admin) error {
	return tx.First(admin, "email=?", admin.Email).Error
}

// FindById implements AdminRepository.
func (repository *AdminRepositoryImpl) FindById(tx *gorm.DB, admin *entity.Admin) error {
	return tx.First(admin).Error
}

// FindAll implements AdminRepository.
func (repository *AdminRepositoryImpl) FindAll(tx *gorm.DB, adminId uint, admins *[]entity.Admin) error {
	return tx.Not("id = ?", adminId).Find(admins).Error
}
