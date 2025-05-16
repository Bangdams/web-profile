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

type ContentUsecase interface {
	Create(ctx context.Context, request *model.ContentCreateRequest) (*model.ContentResponse, error)
	Update(ctx context.Context, request *model.ContentUpdateRequest) (*model.ContentResponse, error)
	Delete(ctx context.Context, contentId uint) error
	FindAll(ctx context.Context, order string, category string) (*[]model.ContentResponse, error)
	FindById(ctx context.Context, contentId uint) (*model.ContentResponse, error)
}

type ContentUsecaseImpl struct {
	ContentRepo repository.ContentRepository
	AdminRepo   repository.AdminRepository
	DB          *gorm.DB
	Validate    *validator.Validate
}

func NewContentUsecase(contentRepo repository.ContentRepository, adminRepo repository.AdminRepository, DB *gorm.DB, validate *validator.Validate) ContentUsecase {
	return &ContentUsecaseImpl{
		ContentRepo: contentRepo,
		AdminRepo:   adminRepo,
		DB:          DB,
		Validate:    validate,
	}
}

// Create implements ContentUsecase.
func (contentUsecase *ContentUsecaseImpl) Create(ctx context.Context, request *model.ContentCreateRequest) (*model.ContentResponse, error) {
	tx := contentUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := contentUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error create content : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	content := &entity.Content{
		Title:       request.Title,
		Description: request.Description,
		Image:       request.Image,
		Address:     request.Address,
		ContactInfo: request.ContactInfo,
		Category:    request.Category,
		CreatedBy:   request.CreatedBy,
	}

	if err := contentUsecase.ContentRepo.Create(tx, content); err != nil {
		log.Println("failed when create repo content : ", err)
		return nil, fiber.ErrInternalServerError
	}

	admin := &entity.Admin{
		ID: request.CreatedBy,
	}

	if err := contentUsecase.AdminRepo.FindById(tx, admin); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Admin data was not found",
				Details: []string{},
			}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id content usecase : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		} else {
			log.Println("Error find by id content usecase:", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	content.Admin = *admin

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase content")
	return converter.ContentToResponse(content), nil
}

// Delete implements ContentUsecase.
func (contentUsecase *ContentUsecaseImpl) Delete(ctx context.Context, contentId uint) error {
	tx := contentUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	content := &entity.Content{
		ID: contentId,
	}

	err := contentUsecase.ContentRepo.FindById(tx, content)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "content data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error delete content : ", err)

			return fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error delete content : ", err)
		return fiber.ErrInternalServerError
	}

	err = contentUsecase.ContentRepo.Delete(tx, content)
	if err != nil {
		log.Println("failed when delete repo content : ", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	log.Println("success delete from usecase content")

	return nil
}

// FindAll implements ContentUsecase.
func (contentUsecase *ContentUsecaseImpl) FindAll(ctx context.Context, order string, category string) (*[]model.ContentResponse, error) {
	tx := contentUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var contents = &[]entity.Content{}
	category = strings.ToLower(category)

	if category != "wisata" && category != "kuliner" && category != "kerajinan" {
		category = ""
	}

	err := contentUsecase.ContentRepo.FindAll(tx, strings.ToUpper(order), category, contents)
	if err != nil {
		log.Println("failed when find all repo content : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find all from usecase content")
	return converter.ContentToResponses(contents), nil
}

// FindById implements ContentUsecase.
func (contentUsecase *ContentUsecaseImpl) FindById(ctx context.Context, contentId uint) (*model.ContentResponse, error) {
	tx := contentUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	content := new(entity.Content)
	content.ID = contentId

	if err := contentUsecase.ContentRepo.FindById(tx, content); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Content data was not found",
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
	return converter.ContentToResponse(content), nil
}

// Update implements ContentUsecase.
func (contentUsecase *ContentUsecaseImpl) Update(ctx context.Context, request *model.ContentUpdateRequest) (*model.ContentResponse, error) {
	tx := contentUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := contentUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error update content : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	content := &entity.Content{
		ID:          request.ID,
		Title:       request.Title,
		Description: request.Description,
		Image:       request.Image,
		Address:     request.Address,
		ContactInfo: request.ContactInfo,
		Category:    request.Category,
		CreatedBy:   request.CreatedBy,
	}

	if err := contentUsecase.ContentRepo.Update(tx, content); err != nil {
		log.Println("failed when update repo content : ", err)
		return nil, fiber.ErrInternalServerError
	}

	admin := &entity.Admin{
		ID: request.CreatedBy,
	}

	if err := contentUsecase.AdminRepo.FindById(tx, admin); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Admin data was not found",
				Details: []string{},
			}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id content usecase : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		} else {
			log.Println("Error find by id content usecase:", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	content.Admin = *admin

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success update from usecase content")
	return converter.ContentToResponse(content), nil
}
