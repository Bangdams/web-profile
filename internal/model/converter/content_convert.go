package converter

import (
	"log"

	"github.com/Bangdams/web-profile-API/internal/entity"
	"github.com/Bangdams/web-profile-API/internal/model"
)

func ContentToResponse(content *entity.Content) *model.ContentResponse {
	log.Println("log from content to response")

	return &model.ContentResponse{
		ID:          content.ID,
		Title:       content.Title,
		Content:     content.Content,
		Image:       content.Image,
		Address:     content.Address,
		ContactInfo: content.ContactInfo,
		Category:    content.Category,
		CreatedBy:   content.Admin.Name,
		CreatedAt:   content.CreatedAt.Format("2006-01-02"),
	}
}

func ContentToResponses(contents *[]entity.Content) *[]model.ContentResponse {
	var contentResponses []model.ContentResponse

	log.Println("log from content to responses")

	for _, content := range *contents {
		contentResponses = append(contentResponses, *ContentToResponse(&content))
	}

	return &contentResponses
}
