package converter

import (
	"log"

	"github.com/Bangdams/web-profile-API/internal/entity"
	"github.com/Bangdams/web-profile-API/internal/model"
)

func AnnouncementToResponse(announcement *entity.Announcement) *model.AnnouncementResponse {
	log.Println("log from announcement to response")

	return &model.AnnouncementResponse{
		ID:          announcement.ID,
		Title:       announcement.Title,
		Content:     announcement.Content,
		Image:       announcement.Image,
		PublishedBy: announcement.Admin.Name,
		CreatedAt:   announcement.CreatedAt.Format("2006-01-02"),
	}
}

func AnnouncementToResponses(announcements *[]entity.Announcement) *[]model.AnnouncementResponse {
	var announcementResponses []model.AnnouncementResponse

	log.Println("log from announcement to responses")

	for _, announcement := range *announcements {
		announcementResponses = append(announcementResponses, *AnnouncementToResponse(&announcement))
	}

	return &announcementResponses
}
