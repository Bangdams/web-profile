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

type ContentController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindById(ctx *fiber.Ctx) error
}

type ContentControllerImpl struct {
	ContentUsecase usecase.ContentUsecase
}

func NewContentController(ContentUsecase usecase.ContentUsecase) ContentController {
	return &ContentControllerImpl{
		ContentUsecase: ContentUsecase,
	}
}

// FindById implements ContentController.
func (controller *ContentControllerImpl) FindById(ctx *fiber.Ctx) error {
	contentId, err := ctx.ParamsInt("content_id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	response, err := controller.ContentUsecase.FindById(ctx.UserContext(), uint(contentId))
	if err != nil {
		log.Println("failed to find by id content")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContentResponse]{Data: response})
}

// Create implements ContentController.
func (controller *ContentControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(model.ContentCreateRequest)

	createdBy, err := strconv.Atoi(ctx.FormValue("created_by"))
	if err != nil {
		log.Println("error badrequest")
		return fiber.ErrBadRequest
	}

	request.Title = ctx.FormValue("title")
	request.Description = ctx.FormValue("description")
	request.Address = ctx.FormValue("address")
	request.ContactInfo = ctx.FormValue("contact_info")
	request.Category = ctx.FormValue("category")
	request.CreatedBy = uint(createdBy)

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

	response, err := controller.ContentUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create content")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContentResponse]{Data: response})
}

// Delete implements ContentController.
func (controller *ContentControllerImpl) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := controller.ContentUsecase.Delete(ctx.UserContext(), uint(id)); err != nil {
		log.Println("failed to delete content")
		return err
	}

	return nil
}

// FindAll implements ContentController.
func (controller *ContentControllerImpl) FindAll(ctx *fiber.Ctx) error {
	var responses *[]model.ContentResponse
	var err error

	order := ctx.Query("order")
	category := ctx.Query("category")
	responses, err = controller.ContentUsecase.FindAll(ctx.UserContext(), order, category)
	if err != nil {
		log.Println("failed to find all content")
		return err
	}

	return ctx.JSON(model.WebResponses[model.ContentResponse]{Data: responses})
}

// Update implements ContentController.
func (controller *ContentControllerImpl) Update(ctx *fiber.Ctx) error {
	request := new(model.ContentUpdateRequest)

	createdBy, err := strconv.Atoi(ctx.FormValue("created_by"))
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
	request.Description = ctx.FormValue("description")
	request.Address = ctx.FormValue("address")
	request.ContactInfo = ctx.FormValue("contact_info")
	request.Category = ctx.FormValue("category")
	request.CreatedBy = uint(createdBy)

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

	response, err := controller.ContentUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create content")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContentResponse]{Data: response})
}
