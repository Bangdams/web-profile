package config

import (
	"github.com/Bangdams/web-profile-API/internal/delivery/http"
	"github.com/Bangdams/web-profile-API/internal/delivery/http/route"
	"github.com/Bangdams/web-profile-API/internal/repository"
	"github.com/Bangdams/web-profile-API/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	// repo
	adminRepo := repository.NewAdminRepository()
	refreshTokenRepo := repository.NewRefreshTokenRepository()

	// usecase
	adminUsecase := usecase.NewAdminUsecase(adminRepo, refreshTokenRepo, config.DB, config.Validate)

	// controller
	adminController := http.NewAdminController(adminUsecase)

	routeConfig := route.RouteConfig{
		App:             config.App,
		AdminController: adminController,
	}

	routeConfig.Setup()
}
