package model

type AnnouncementResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Image       string `json:"image"`
	PublishedBy string `json:"published_by"`
	CreatedAt   string `json:"created_at"`
}

type AnnouncementCreateRequest struct {
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content" validate:"required"`
	Image       string `json:"image" validate:"required"`
	PublishedBy uint   `json:"published_by" validate:"required"`
}

type AnnouncementUpdateRequest struct {
	ID uint `json:"id" validate:"required"`
	AnnouncementCreateRequest
}
