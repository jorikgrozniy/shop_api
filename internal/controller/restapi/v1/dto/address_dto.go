package dto

type AddressRequest struct {
	Country string `json:"country" binding:"required" example:"Russia"`
	City    string `json:"city" binding:"required" example:"Moscow"`
	Street  string `json:"street" binding:"required" example:"Vernadskogo"`
}

type AddressResponse struct {
	ID      string `json:"address_id" example:"b1c4b160-b9e9-4f0c-83ec-9119e1f467f9"`
	Country string `json:"country" example:"Russia"`
	City    string `json:"city" example:"Moscow"`
	Street  string `json:"street" example:"Vernadskogo"`
}
