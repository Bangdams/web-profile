package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Bangdams/web-profile-API/internal/entity"
	"github.com/Bangdams/web-profile-API/internal/model"
	"github.com/Bangdams/web-profile-API/internal/model/converter"
	"github.com/Bangdams/web-profile-API/internal/repository"
	"github.com/Bangdams/web-profile-API/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminUsecase interface {
	Create(ctx context.Context, request *model.AdminCreateRequest) (*model.AdminResponse, error)
	Update(ctx context.Context, request *model.AdminUpdateRequest) (*model.AdminResponse, error)
	Delete(ctx context.Context, adminId uint) error
	FindAll(ctx context.Context, adminId uint) (*[]model.AdminResponse, error)
	FindByUsername(ctx context.Context, usernameRequest string) (*model.AdminResponse, error)
	Login(ctx context.Context, request *model.LoginRequest, requestRefreshToken string) (*model.LoginResponse, string, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, error)
}

type AdminUsecaseImpl struct {
	AdminRepo        repository.AdminRepository
	RefreshTokenRepo repository.RefreshTokenRepository
	DB               *gorm.DB
	Validate         *validator.Validate
}

func NewAdminUsecase(adminRepo repository.AdminRepository, refreshTokenRepo repository.RefreshTokenRepository, DB *gorm.DB, validate *validator.Validate) AdminUsecase {
	return &AdminUsecaseImpl{
		AdminRepo:        adminRepo,
		RefreshTokenRepo: refreshTokenRepo,
		DB:               DB,
		Validate:         validate,
	}
}

// Login implements AdminUsecase.
func (adminUsecase *AdminUsecaseImpl) Login(ctx context.Context, request *model.LoginRequest, requestRefreshTokenAdmin string) (*model.LoginResponse, string, error) {
	now := time.Now()

	_, err := util.ParseToken(requestRefreshTokenAdmin, []byte(os.Getenv("SECRET_KEY")))
	if err == nil {
		return nil, "", fiber.NewError(fiber.StatusBadRequest, "Refresh token still valid")
	}

	tx := adminUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	admin := &entity.Admin{}

	if err := adminUsecase.AdminRepo.Login(tx, admin, request.Username); err != nil {
		log.Println("invalid Username : ", err)
		return nil, "", fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.Password)); err != nil {
		log.Println("invalid password : ", err)
		return nil, "", fiber.ErrUnauthorized
	}

	accessToken, err := util.GenerateAccessToken(admin)
	if err != nil {
		log.Println("Failed to generate token jwt")
		return nil, "", fiber.ErrInternalServerError
	}

	refreshToken, err := util.GenerateRefreshToken(admin)
	if err != nil {
		log.Println("Failed to generate token jwt")
		return nil, "", fiber.ErrInternalServerError
	}

	duration := os.Getenv("DURATION_JWT_REFRESH_TOKEN")
	lifeTime, _ := strconv.Atoi(duration)

	requestRefreshToken := &entity.RefreshToken{
		AdminId:      admin.ID,
		StatusLogout: 0,
		Token:        refreshToken,
		ExpiresAt:    now.Add(time.Minute * time.Duration(lifeTime)),
	}

	if err := adminUsecase.RefreshTokenRepo.FindById(tx, admin.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Creating new refresh token in database")
			adminUsecase.RefreshTokenRepo.Create(tx, requestRefreshToken)
		} else {
			log.Println("Error fetching refresh token:", err)
			return nil, "", fiber.ErrInternalServerError
		}
	} else {
		log.Println("Updating existing refresh token in database")
		adminUsecase.RefreshTokenRepo.Update(tx, requestRefreshToken)
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, "", fiber.ErrInternalServerError
	}

	log.Println("success login")

	return converter.LoginAdminToResponse(accessToken), refreshToken, nil
}

// Logout implements AdminUsecase.
func (adminUsecase *AdminUsecaseImpl) Logout(ctx context.Context, refreshToken string) error {
	tx := adminUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	claims, err := util.ParseToken(refreshToken, []byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return fiber.ErrUnauthorized
	}

	adminId := claims["admin_id"].(float64)

	if err := adminUsecase.RefreshTokenRepo.FindById(tx, uint(adminId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Logout failed because the user has not logged in.")
			return fiber.NewError(fiber.StatusBadRequest, "User has not logged in")
		} else {
			log.Println("Error RefreshToken findbyid:", err)
			return fiber.ErrInternalServerError
		}
	} else {
		log.Println("Logout successful.")
		adminUsecase.RefreshTokenRepo.Update(tx, &entity.RefreshToken{AdminId: uint(adminId), StatusLogout: 1})
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	return nil
}

// Refresh implements AdminUsecase.
func (adminUsecase *AdminUsecaseImpl) Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {
	tx := adminUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	claims, err := util.ParseToken(refreshToken, []byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	adminId := claims["admin_id"].(float64)
	request := entity.Admin{
		ID:       uint(adminId),
		Name:     claims["name"].(string),
		Username: claims["email"].(string),
	}

	if err := adminUsecase.RefreshTokenRepo.CheckStatusLogout(tx, uint(adminId)); err != nil {
		return nil, fiber.ErrUnauthorized
	}

	newAccessToken, _ := util.GenerateAccessToken(&request)

	log.Println("success create access token")

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.LoginAdminToResponse(newAccessToken), nil
}

// Create implements AdminUsecase.
func (adminUsecase *AdminUsecaseImpl) Create(ctx context.Context, request *model.AdminCreateRequest) (*model.AdminResponse, error) {
	tx := adminUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := adminUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error create admin : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to generate password")
		return nil, fiber.ErrInternalServerError
	}

	admin := &entity.Admin{
		Name:     request.Name,
		Username: request.Username,
		Password: string(password),
	}

	if err := adminUsecase.AdminRepo.FindByUsername(tx, admin); err == nil {
		errorResponse.Message = "Duplicate entry"
		errorResponse.Details = []string{"email already exists in the database."}

		jsonString, _ := json.Marshal(errorResponse)

		return nil, fiber.NewError(fiber.ErrConflict.Code, string(jsonString))
	}

	err = adminUsecase.AdminRepo.Create(tx, admin)
	if err != nil {
		log.Println("failed when create repo admin : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase admin")

	return converter.AdminToResponse(admin), nil
}

// Delete implements AdminUsecase.
func (adminUsecase *AdminUsecaseImpl) Delete(ctx context.Context, adminId uint) error {
	tx := adminUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	admin := &entity.Admin{
		ID: adminId,
	}

	err := adminUsecase.AdminRepo.FindById(tx, admin)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "admin data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error delete admin : ", err)

			return fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error delete admin : ", err)
		return fiber.ErrInternalServerError
	}

	err = adminUsecase.AdminRepo.Delete(tx, admin)
	if err != nil {
		log.Println("failed when delete repo admin : ", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	log.Println("success delete from usecase admin")

	return nil
}

// FindAll implements AdminUsecase.
func (adminUsecase *AdminUsecaseImpl) FindAll(ctx context.Context, adminId uint) (*[]model.AdminResponse, error) {
	tx := adminUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var admins = &[]entity.Admin{}
	err := adminUsecase.AdminRepo.FindAll(tx, adminId, admins)
	if err != nil {
		log.Println("failed when find all repo admin : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find all from usecase admin")

	return converter.AdminToResponses(admins), nil
}

// FindById implements AdminUsecase.
func (adminUsecase *AdminUsecaseImpl) FindByUsername(ctx context.Context, emailRequest string) (*model.AdminResponse, error) {
	tx := adminUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	admin := new(entity.Admin)
	admin.Username = emailRequest

	if err := adminUsecase.AdminRepo.FindByUsername(tx, admin); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Admin data was not found",
				Details: []string{},
			}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by email admin usecase : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		} else {
			log.Println("Error find by email admin usecase:", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find by email from usecase admin")

	return converter.AdminToResponse(admin), nil
}

// Update implements AdminUsecase.
func (adminUsecase *AdminUsecaseImpl) Update(ctx context.Context, request *model.AdminUpdateRequest) (*model.AdminResponse, error) {
	tx := adminUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := adminUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error update admin : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	admin := &entity.Admin{
		Username: request.Username,
	}

	err = adminUsecase.AdminRepo.FindByUsername(tx, admin)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "Admin data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)
			log.Println("Data not found")

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by email : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("failed to generate password")
			return nil, fiber.ErrInternalServerError
		}

		admin.Password = string(password)
	}

	admin.Name = request.Name

	err = adminUsecase.AdminRepo.Update(tx, admin)
	if err != nil {
		mysqlErr := err.(*mysql.MySQLError)
		log.Println("failed when update repo admin : ", err)

		var errorField string
		parts := strings.Split(mysqlErr.Message, "'")
		if len(parts) > 2 {
			errorField = parts[1]
		}

		if mysqlErr.Number == 1062 {
			errorResponse.Message = "Duplicate entry"
			errorResponse.Details = []string{errorField + " already exists in the database."}

			jsonString, _ := json.Marshal(errorResponse)

			return nil, fiber.NewError(fiber.ErrConflict.Code, string(jsonString))
		}

		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success update from usecase admin")
	return converter.AdminToResponse(admin), nil
}
