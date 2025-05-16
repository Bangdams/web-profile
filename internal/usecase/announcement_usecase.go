package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Bangdams/web-profile-API/internal/entity"
	"github.com/Bangdams/web-profile-API/internal/model"
	"github.com/Bangdams/web-profile-API/internal/model/converter"
	"github.com/Bangdams/web-profile-API/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AnnouncementUsecase interface {
	Create(ctx context.Context, request *model.AnnouncementCreateRequest) (*model.AnnouncementResponse, error)
	Update(ctx context.Context, request *model.AnnouncementUpdateRequest) (*model.AnnouncementResponse, error)
	Delete(ctx context.Context, announcemenId uint) error
	FindAll(ctx context.Context, order string) (*[]model.AnnouncementResponse, error)
	FindById(ctx context.Context, announcementtId uint) (*model.AnnouncementResponse, error)
}

type AnnouncementUsecaseImpl struct {
	AnnouncementRepo repository.AnnouncementRepository
	AdminRepo        repository.AdminRepository
	DB               *gorm.DB
	Validate         *validator.Validate
}

func NewAnnouncementUsecase(announcementRepo repository.AnnouncementRepository, adminRepo repository.AdminRepository, DB *gorm.DB, validate *validator.Validate) AnnouncementUsecase {
	return &AnnouncementUsecaseImpl{
		AnnouncementRepo: announcementRepo,
		AdminRepo:        adminRepo,
		DB:               DB,
		Validate:         validate,
	}
}

// Create implements AnnouncementUsecase.
func (announcementUsecase *AnnouncementUsecaseImpl) Create(ctx context.Context, request *model.AnnouncementCreateRequest) (*model.AnnouncementResponse, error) {
	tx := announcementUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := announcementUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error create announcement : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	announcement := entity.Announcement{
		Title:       request.Title,
		Content:     request.Content,
		Image:       request.Image,
		PublishedBy: request.PublishedBy,
	}

	if err := announcementUsecase.AnnouncementRepo.Create(tx, &announcement); err != nil {
		log.Println("failed when create repo announcement : ", err)
		return nil, fiber.ErrInternalServerError
	}

	admin := &entity.Admin{
		ID: request.PublishedBy,
	}

	if err := announcementUsecase.AdminRepo.FindById(tx, admin); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Admin data was not found",
				Details: []string{},
			}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id announcement usecase : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		} else {
			log.Println("Error find by id announcement usecase:", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	announcement.Admin = *admin
	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase announcement")
	return converter.AnnouncementToResponse(&announcement), nil

}

// Delete implements AnnouncementUsecase.
func (announcementUsecase *AnnouncementUsecaseImpl) Delete(ctx context.Context, announcementId uint) error {
	tx := announcementUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	announcement := &entity.Announcement{
		ID: announcementId,
	}

	err := announcementUsecase.AnnouncementRepo.FindById(tx, announcement)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "announcement data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error delete announcement : ", err)

			return fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error delete announcement : ", err)
		return fiber.ErrInternalServerError
	}

	err = announcementUsecase.AnnouncementRepo.Delete(tx, announcement)
	if err != nil {
		log.Println("failed when delete repo announcement : ", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	log.Println("success delete from usecase announcement")

	return nil
}

// FindAll implements AnnouncementUsecase.
func (announcementUsecase *AnnouncementUsecaseImpl) FindAll(ctx context.Context, order string) (*[]model.AnnouncementResponse, error) {
	tx := announcementUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var announcements = &[]entity.Announcement{}
	err := announcementUsecase.AnnouncementRepo.FindAll(tx, strings.ToUpper(order), announcements)
	if err != nil {
		log.Println("failed when find all repo announcement : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find all from usecase announcement")
	return converter.AnnouncementToResponses(announcements), nil
}

// FindById implements AnnouncementUsecase.
func (announcementUsecase *AnnouncementUsecaseImpl) FindById(ctx context.Context, announcementId uint) (*model.AnnouncementResponse, error) {
	tx := announcementUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	announcement := new(entity.Announcement)
	announcement.ID = announcementId

	if err := announcementUsecase.AnnouncementRepo.FindById(tx, announcement); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Announcement data was not found",
				Details: []string{},
			}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id conten usecase : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		} else {
			log.Println("Error find by id conten usecase:", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}
	return converter.AnnouncementToResponse(announcement), nil
}

// Update implements AnnouncementUsecase.
func (announcementUsecase *AnnouncementUsecaseImpl) Update(ctx context.Context, request *model.AnnouncementUpdateRequest) (*model.AnnouncementResponse, error) {
	tx := announcementUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := announcementUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error update announcement : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	announcement := entity.Announcement{
		ID:          request.ID,
		Title:       request.Title,
		Content:     request.Content,
		Image:       request.Image,
		PublishedBy: request.PublishedBy,
	}

	if err := announcementUsecase.AnnouncementRepo.Update(tx, &announcement); err != nil {
		log.Println("failed when update repo announcement : ", err)
		return nil, fiber.ErrInternalServerError
	}

	admin := &entity.Admin{
		ID: request.PublishedBy,
	}

	if err := announcementUsecase.AdminRepo.FindById(tx, admin); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Admin data was not found",
				Details: []string{},
			}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id announcement usecase : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		} else {
			log.Println("Error find by id announcement usecase:", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	announcement.Admin = *admin

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success update from usecase announcement")
	return converter.AnnouncementToResponse(&announcement), nil

}
