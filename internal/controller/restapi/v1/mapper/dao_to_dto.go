package mapper

import (
	"shop_api/internal/controller/restapi/v1/dto"
	"shop_api/internal/entity"
)

func AddressDAOtoDTO(address *entity.Address) dto.AddressResponse {
	return dto.AddressResponse{
		ID:      address.ID.String(),
		Country: address.Country,
		City:    address.City,
		Street:  address.Street,
	}
}

func ImageDAOtoDTO(image *entity.Image) dto.ImageRequest {
	return dto.ImageRequest{
		JPGbase64: string(image.Image),
	}
}

func ClientDAOtoDTO(client *entity.Client) dto.ClientResponse {
	return dto.ClientResponse{
		ID:               client.ID.String(),
		Name:             client.Name,
		Surname:          client.Surname,
		Birthdate:        client.Birthdate.Format("2006-01-02"),
		Gender:           client.Gender,
		RegistrationDate: client.RegistrationDate.Format("2006-01-02"),
	}
}

func ClientListDAOtoDTO(clients []*entity.Client) dto.ClientListResponse {
	arrLen := len(clients)
	dtoClients := make([]dto.ClientResponse, arrLen)

	for i, client := range clients {
		dtoClients[i] = ClientDAOtoDTO(client)
	}

	return dto.ClientListResponse{
		Data:  dtoClients,
		Total: arrLen,
	}
}

func SupplierDAOtoDTO(supplier *entity.Supplier) dto.SupplierResponse {
	return dto.SupplierResponse{
		ID:          supplier.ID.String(),
		Name:        supplier.Name,
		PhoneNumber: supplier.PhoneNumber,
	}
}

func SupplierListDAOtoDTO(suppliers []*entity.Supplier) dto.SupplierListResponse {
	arrLen := len(suppliers)
	dtoSuppliers := make([]dto.SupplierResponse, arrLen)

	for i, supplier := range suppliers {
		dtoSuppliers[i] = SupplierDAOtoDTO(supplier)
	}

	return dto.SupplierListResponse{
		Data:  dtoSuppliers,
		Total: arrLen,
	}
}

func ProductDAOtoDTO(product *entity.Product) dto.ProductResponse {
	imageID := ""
	if product.ImageID != nil {
		imageID = product.ImageID.String()
	}

	supplierID := ""
	if product.SupplierID != nil {
		supplierID = product.SupplierID.String()
	}

	return dto.ProductResponse{
		ID:             product.ID.String(),
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		LastUpdate:     product.LastUpdate.String(),
		SupplierID:     supplierID,
		ImageID:        imageID,
	}
}

func ProductListDAOtoDTO(products []*entity.Product) dto.ProductListResponse {
	arrLen := len(products)
	dtoProducts := make([]dto.ProductResponse, arrLen)

	for i, product := range products {
		dtoProducts[i] = ProductDAOtoDTO(product)
	}

	return dto.ProductListResponse{
		Data:  dtoProducts,
		Total: arrLen,
	}
}
