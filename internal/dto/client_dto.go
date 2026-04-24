package dto

type ClientRequest struct {
	Name      string         `json:"name" binding:"required"`
	Surname   string         `json:"surname" binding:"required"`
	Birthdate string         `json:"birthdate" binding:"required"`
	Gender    string         `json:"country" binding:"required,oneof=male female"`
	Address   AddressRequest `json:"address" binding:"required"`
}

type SearchClientRequest struct {
	Name    string `form:"name" binding:"required"`
	Surname string `form:"surname" binding:"required"`
}

type ClientResponse struct {
	ID               string          `json:"client_id"`
	Name             string          `json:"name"`
	Surname          string          `json:"surname"`
	Birthdate        string          `json:"birthdate"`
	Gender           string          `json:"country"`
	RegistrationDate string          `json:"registration_date"`
	Address          AddressResponse `json:"address"`
}

type ClientListResponse struct {
	Data   []ClientResponse `json:"data"`
	Total  int              `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}
