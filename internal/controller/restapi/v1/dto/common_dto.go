package dto

type ListRequest struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100" example:"20"`
	Offset int `form:"offset" binding:"omitempty,min=0" example:"0"`
}
