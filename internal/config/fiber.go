package config

import (
	"encoding/json"

	middelware "github.com/Bangdams/web-profile-API/internal/delivery/http/middleware"
	"github.com/Bangdams/web-profile-API/internal/model"
	"github.com/gofiber/fiber/v2"
)

func NewFiber() *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      "Quizku API",
		ErrorHandler: NewErrorHandler(),
	})

	middelware.Middelware(app)

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		var errorResponse model.ErrorResponse
		jsonError := json.Unmarshal([]byte(err.Error()), &errorResponse)
		if jsonError != nil {
			errorResponse.Message = err.Error()
			return ctx.Status(code).JSON(model.WebResponse[any]{Errors: &errorResponse})
		}

		return ctx.Status(code).JSON(model.WebResponse[any]{Errors: &errorResponse})
	}
}
