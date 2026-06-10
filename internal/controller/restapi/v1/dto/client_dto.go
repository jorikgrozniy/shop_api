package dto

type ClientRequest struct {
	Name      string         `json:"name" binding:"required" example:"Egor"`
	Surname   string         `json:"surname" binding:"required" example:"Egorov"`
	Birthdate string         `json:"birthdate" binding:"required" example:"2006-11-25"`
	Gender    string         `json:"gender" binding:"required,oneof=male female" example:"male"`
	Address   AddressRequest `json:"address" binding:"required"`
}

type ClientListRequest struct {
	Name    string `form:"name" binding:"omitempty" example:"Egor"`
	Surname string `form:"surname" binding:"omitempty" example:"Egorov"`
	Limit   int    `form:"limit" binding:"omitempty,min=1,max=100" example:"20"`
	Offset  int    `form:"offset" binding:"omitempty,min=0" example:"0"`
}

type ClientResponse struct {
	ID               string          `json:"client_id" example:"b1c4b160-b9e9-4f0c-83ec-9119e1f467f9"`
	Name             string          `json:"name" example:"Egor"`
	Surname          string          `json:"surname" example:"Egorov"`
	Birthdate        string          `json:"birthdate" example:"2006-11-25"`
	Gender           string          `json:"country" example:"male"`
	RegistrationDate string          `json:"registration_date" example:"2026-06-01"`
	Address          AddressResponse `json:"address"`
}

type ClientListResponse struct {
	Data   []ClientResponse `json:"data"`
	Total  int              `json:"total" example:"5"`
	Limit  int              `json:"limit" example:"20"`
	Offset int              `json:"offset" example:"0"`
}
