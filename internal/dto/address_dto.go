package dto

type AddressRequest struct {
	Country string `json:"country" binding:"required"`
	City    string `json:"city" binding:"required"`
	Street  string `json:"street" binding:"required"`
}

type AddressResponse struct {
	ID      string `json:"address_id"`
	Country string `json:"country"`
	City    string `json:"city"`
	Street  string `json:"street"`
}
