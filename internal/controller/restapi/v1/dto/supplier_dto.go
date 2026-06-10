package dto

type SupplierRequest struct {
	Name        string         `json:"name" binding:"required" example:"Super Keyboard Supplier"`
	PhoneNumber string         `json:"phone_number" binding:"required" example:"+79995351675"`
	Address     AddressRequest `json:"address" binding:"required"`
}

type SupplierResponse struct {
	ID          string          `json:"supplier_id" example:"b1c4b160-b9e9-4f0c-83ec-9119e1f467f9"`
	Name        string          `json:"name" example:"Super Keyboard Supplier"`
	PhoneNumber string          `json:"phone_number" example:"+79995351675"`
	Address     AddressResponse `json:"address"`
}

type SupplierListResponse struct {
	Data   []SupplierResponse `json:"data"`
	Total  int                `json:"total" example:"5"`
	Limit  int                `json:"limit" example:"20"`
	Offset int                `json:"offset" example:"0"`
}
