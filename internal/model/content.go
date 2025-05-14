package model

type ContentResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Address     string `json:"address"`
	ContactInfo string `json:"contact_info"`
	Category    string `json:"category"`
	CreatedBy   string `json:"created_by"`
}

type ContentCreateRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Image       string `json:"image" validate:"required"`
	Address     string `json:"address" validate:"required"`
	ContactInfo string `json:"contact_info" validate:"required"`
	Category    string `json:"category" validate:"required"`
	CreatedBy   uint   `json:"created_by" validate:"required"`
}

type ContentUpdateRequest struct {
	ID uint `json:"id" validate:"required"`
	ContentCreateRequest
}
