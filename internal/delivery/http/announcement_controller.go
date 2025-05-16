package http

import (
	"log"
	"path/filepath"
	"strconv"

	"github.com/Bangdams/web-profile-API/internal/model"
	"github.com/Bangdams/web-profile-API/internal/usecase"
	"github.com/Bangdams/web-profile-API/internal/util"
	"github.com/gofiber/fiber/v2"
)

type AnnouncementController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindById(ctx *fiber.Ctx) error
}

type AnnouncementControllerImpl struct {
	AnnouncementUsecase usecase.AnnouncementUsecase
}

func NewAnnouncementController(AnnouncementUsecase usecase.AnnouncementUsecase) AnnouncementController {
	return &AnnouncementControllerImpl{
		AnnouncementUsecase: AnnouncementUsecase,
	}
}

// Create implements AnnouncementController.
func (controller *AnnouncementControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(model.AnnouncementCreateRequest)

	publishedBy, err := strconv.Atoi(ctx.FormValue("published_by"))
	if err != nil {
		log.Println("error badrequest")
		return fiber.ErrBadRequest
	}

	request.Title = ctx.FormValue("title")
	request.Content = ctx.FormValue("content")
	request.PublishedBy = uint(publishedBy)

	// upload image
	file, err := ctx.FormFile("image")
	if err != nil {
		log.Println("failed to parse request image : ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "image is required"})
	}

	filename := filepath.Base(file.Filename)
	generateFilename := util.GenerateRandomFilename(filename)
	savePath := filepath.Join("./upload", generateFilename)

	if err := ctx.SaveFile(file, savePath); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save image"})
	}

	request.Image = generateFilename
	// end upload image

	response, err := controller.AnnouncementUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create announcement")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AnnouncementResponse]{Data: response})
}

// Delete implements AnnouncementController.
func (controller *AnnouncementControllerImpl) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := controller.AnnouncementUsecase.Delete(ctx.UserContext(), uint(id)); err != nil {
		log.Println("failed to delete announcement")
		return err
	}

	return nil
}

// FindAll implements AnnouncementController.
func (controller *AnnouncementControllerImpl) FindAll(ctx *fiber.Ctx) error {
	var responses *[]model.AnnouncementResponse
	var err error

	order := ctx.Query("order")
	responses, err = controller.AnnouncementUsecase.FindAll(ctx.UserContext(), order)
	if err != nil {
		log.Println("failed to find all announcement")
		return err
	}

	return ctx.JSON(model.WebResponses[model.AnnouncementResponse]{Data: responses})
}

// FindById implements AnnouncementController.
func (controller *AnnouncementControllerImpl) FindById(ctx *fiber.Ctx) error {
	announcementId, err := ctx.ParamsInt("announcement_id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	response, err := controller.AnnouncementUsecase.FindById(ctx.UserContext(), uint(announcementId))
	if err != nil {
		log.Println("failed to find by id announcement")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AnnouncementResponse]{Data: response})
}

// Update implements AnnouncementController.
func (controller *AnnouncementControllerImpl) Update(ctx *fiber.Ctx) error {
	request := new(model.AnnouncementUpdateRequest)

	publishedBy, err := strconv.Atoi(ctx.FormValue("published_by"))
	if err != nil {
		log.Println("error bad request : ", err)
		return fiber.ErrBadRequest
	}

	id, err := strconv.Atoi(ctx.FormValue("id"))
	if err != nil {
		log.Println("error bad request : ", err)
		return fiber.ErrBadRequest
	}

	request.ID = uint(id)
	request.Title = ctx.FormValue("title")
	request.Content = ctx.FormValue("content")
	request.PublishedBy = uint(publishedBy)

	// upload image
	var filename string
	filename = ctx.FormValue("image_name")
	file, err := ctx.FormFile("image")
	if err == nil {
		filename = filepath.Base(file.Filename)
		generateFilename := util.GenerateRandomFilename(filename)
		savePath := filepath.Join("./upload", generateFilename)

		if err := ctx.SaveFile(file, savePath); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save image"})
		}

		request.Image = generateFilename
	} else {
		request.Image = filename
	}
	// end upload image

	response, err := controller.AnnouncementUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create announcement")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AnnouncementResponse]{Data: response})
}
