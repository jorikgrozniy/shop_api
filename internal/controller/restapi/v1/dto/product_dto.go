package dto

type ProductRequest struct {
	Name           string       `json:"name" binding:"required" example:"Keyboard"`
	Category       string       `json:"category" binding:"required" example:"Electronics"`
	Price          float64      `json:"price" binding:"required,min=0" example:"289.90"`
	AvailableStock int          `json:"available_stock" binding:"required,min=1" example:"30"`
	SupplierID     string       `json:"supplier_id" binding:"required" example:"b1c4b160-b9e9-4f0c-83ec-9119e1f467f9"`
	Image          ImageRequest `json:"image" binding:"omitempty"`
}

type ProductResponse struct {
	ID             string  `json:"product_id" example:"b1c4b160-b9e9-4f0c-83ec-9119e1f467f9"`
	Name           string  `json:"name" example:"Keyboard"`
	Category       string  `json:"category" example:"Electronics"`
	Price          float64 `json:"price" example:"289.90"`
	AvailableStock int     `json:"available_stock" example:"30"`
	LastUpdate     string  `json:"last_update" example:"2026-05-20"`
	SupplierID     string  `json:"supplier_id,omitempty" example:"b1c4b160-b9e9-4f0c-83ec-9119e1f467f9"`
	ImageID        string  `json:"image_id,omitempty" example:"b1c4b160-b9e9-4f0c-83ec-9119e1f467f9"`
}

type ProductListResponse struct {
	Data   []ProductResponse `json:"data"`
	Total  int               `json:"total" example:"5"`
	Limit  int               `json:"limit" example:"20"`
	Offset int               `json:"offset" example:"0"`
}

type ProductStockRequest struct {
	Amount int `json:"amount" binding:"required,min=1" example:"2"`
}
