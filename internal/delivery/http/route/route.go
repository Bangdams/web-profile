package route

import (
	"github.com/Bangdams/web-profile-API/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App               *fiber.App
	AdminController   http.AdminController
	ContentController http.ContentController
}

func (config *RouteConfig) Setup() {
	// Api for login
	config.App.Post("/login", config.AdminController.Login)
	config.App.Post("/logout", config.AdminController.Logout)
	config.App.Post("/refresh", config.AdminController.Refresh)
	config.App.Get("/api/status-login", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"message": "success"})
	})

	// API for admin
	config.App.Get("/api/admins", config.AdminController.FindAll)
	config.App.Get("/api/admins/:email", config.AdminController.FindByEmail)
	config.App.Post("/api/admins", config.AdminController.Create)
	config.App.Delete("/api/admins/:id", config.AdminController.Delete)
	config.App.Put("/api/admins", config.AdminController.Update)

	// API for content
	config.App.Get("/api/contents", config.ContentController.FindAll)
	config.App.Get("/api/contents/:content_id", config.ContentController.FindById)
	config.App.Post("/api/contents", config.ContentController.Create)
	config.App.Delete("/api/contents/:id", config.ContentController.Delete)
	config.App.Put("/api/contents", config.ContentController.Update)
}
