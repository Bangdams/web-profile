package http

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Bangdams/web-profile-API/internal/model"
	"github.com/Bangdams/web-profile-API/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AdminController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindByUsername(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	Refresh(ctx *fiber.Ctx) error
}

type AdminControllerImpl struct {
	AdminUsecase usecase.AdminUsecase
}

func NewAdminController(AdminUsecase usecase.AdminUsecase) AdminController {
	return &AdminControllerImpl{
		AdminUsecase: AdminUsecase,
	}
}

// Login implements AdminController.
func (controller *AdminControllerImpl) Login(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("refresh_token")

	request := new(model.LoginRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, refreshToken, err := controller.AdminUsecase.Login(ctx.UserContext(), request, cookie)
	if err != nil {
		log.Println("failed to login")
		return err
	}

	// durasi refreshToken
	duration := os.Getenv("DURATION_JWT_REFRESH_TOKEN")
	lifeTime, _ := strconv.Atoi(duration)

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   60 * 60 * 24 * lifeTime,
	})

	return ctx.JSON(model.WebResponse[*model.LoginResponse]{Data: response})
}

// Logout implements AdminController.
func (controller *AdminControllerImpl) Logout(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("refresh_token")
	if cookie == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	err := controller.AdminUsecase.Logout(ctx.UserContext(), cookie)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return ctx.JSON(model.WebResponse[string]{Data: "Logout successful"})
}

// Refresh implements AdminController.
func (controller *AdminControllerImpl) Refresh(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("refresh_token")
	if cookie == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	response, err := controller.AdminUsecase.Refresh(ctx.UserContext(), cookie)
	if err != nil {
		log.Println("failed to create refresh token")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.LoginResponse]{Data: response})
}

// Create implements AdminController.
func (controller *AdminControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(model.AdminCreateRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.AdminUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AdminResponse]{Data: response})
}

// Delete implements AdminController.
func (controller *AdminControllerImpl) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := controller.AdminUsecase.Delete(ctx.UserContext(), uint(id)); err != nil {
		log.Println("failed to delete user")
		return err
	}

	return nil
}

// FindAll implements AdminController.
func (controller *AdminControllerImpl) FindAll(ctx *fiber.Ctx) error {
	var responses *[]model.AdminResponse
	var err error

	adminToken := ctx.Locals("admin").(*jwt.Token)
	claims := adminToken.Claims.(jwt.MapClaims)
	adminId := claims["admin_id"].(float64)

	responses, err = controller.AdminUsecase.FindAll(ctx.UserContext(), uint(adminId))
	if err != nil {
		log.Println("failed to find all admin")
		return err
	}

	return ctx.JSON(model.WebResponses[model.AdminResponse]{Data: responses})
}

// FindByUsername implements AdminController.
func (controller *AdminControllerImpl) FindByUsername(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	response, err := controller.AdminUsecase.FindByUsername(ctx.UserContext(), username)
	if err != nil {
		log.Println("failed to find by username admin")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AdminResponse]{Data: response})
}

// Update implements AdminController.
func (controller *AdminControllerImpl) Update(ctx *fiber.Ctx) error {
	request := new(model.AdminUpdateRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.AdminUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to update admin")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AdminResponse]{Data: response})
}
