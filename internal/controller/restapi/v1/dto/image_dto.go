package dto

type ImageRequest struct {
	JPGbase64 string `json:"jpg_base64" binding:"required" example:"iVBORw0KGgoAAAANSUhEUgAAB4AA..."`
}
