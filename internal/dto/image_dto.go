package dto

type ImageRequest struct {
	Image string `json:"image" binding:"required"`
}

type ImageResponse struct {
	ID    string `json:"image_id"`
	Image string `json:"image"`
}
