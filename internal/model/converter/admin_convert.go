package converter

import (
	"log"

	"github.com/Bangdams/web-profile-API/internal/entity"
	"github.com/Bangdams/web-profile-API/internal/model"
)

func AdminToResponse(admin *entity.Admin) *model.AdminResponse {
	log.Println("log from admin to response")

	return &model.AdminResponse{
		ID:    admin.ID,
		Email: admin.Email,
		Name:  admin.Name,
	}
}

func AdminToResponses(admins *[]entity.Admin) *[]model.AdminResponse {
	var adminResponses []model.AdminResponse

	log.Println("log from admin to responses")

	for _, admin := range *admins {
		adminResponses = append(adminResponses, *AdminToResponse(&admin))
	}

	return &adminResponses
}

func LoginAdminToResponse(token string) *model.LoginResponse {
	log.Println("log from login admin to response")

	return &model.LoginResponse{
		AccessToken: token,
	}
}
