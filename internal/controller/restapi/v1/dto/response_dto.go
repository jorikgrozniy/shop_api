package dto

type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

type SuccessResponse struct {
	Status string `json:"status" example:"created"`
}
