package dto

type SupplierRequest struct {
	Name        string         `json:"name" binding:"required"`
	PhoneNumber string         `json:"phone_number" binding:"required"`
	Address     AddressRequest `json:"address" binding:"required"`
}

type SupplierResponse struct {
	ID          string          `json:"supplier_id"`
	Name        string          `json:"name"`
	PhoneNumber string          `json:"phone_number"`
	Address     AddressResponse `json:"address"`
}

type SupplierListResponse struct {
	Data   []SupplierResponse `json:"data"`
	Total  int                `json:"total"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
}
