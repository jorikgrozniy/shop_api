package dto

type ProductRequest struct {
	Name           string       `json:"name" binding:"required"`
	Category       string       `json:"category" binding:"required"`
	Price          float64      `json:"price" binding:"required,min=0"`
	AvailableStock int          `json:"available_stock" binding:"required,min=1"`
	SupplierID     string       `json:"supplier_id" binding:"required"`
	Image          ImageRequest `json:"image" binding:"omitempty"`
}

type ProductResponse struct {
	ID             string
	Name           string        `json:"name"`
	Category       string        `json:"category"`
	Price          float64       `json:"price"`
	AvailableStock int           `json:"available_stock"`
	LastUpdate     string        `json:"last_update"`
	SupplierID     string        `json:"supplier_id"`
	Image          ImageResponse `json:"image,omitzero"`
}

type ProductListResponse struct {
	Data   []ProductResponse `json:"data"`
	Total  int               `json:"total"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}
