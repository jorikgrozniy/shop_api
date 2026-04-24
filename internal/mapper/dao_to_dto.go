package mapper

import (
	"shop_api/internal/dao"
	"shop_api/internal/dto"
)

func AddressDAOtoDTO(address *dao.Address) dto.AddressResponse {
	return dto.AddressResponse{
		ID:      address.ID.String(),
		Country: address.Country,
		City:    address.City,
		Street:  address.Street,
	}
}

func ImageDAOtoDTO(image *dao.Image) dto.ImageResponse {
	return dto.ImageResponse{
		ID:    image.ID.String(),
		Image: string(image.Image),
	}
}

func ClientDAOtoDTO(client *dao.Client) dto.ClientResponse {
	return dto.ClientResponse{
		ID:               client.ID.String(),
		Name:             client.Name,
		Surname:          client.Surname,
		Birthdate:        client.Birthdate.String(),
		Gender:           client.Gender,
		RegistrationDate: client.RegistrationDate.String(),
	}
}

func ClientListDAOtoDTO(clients []*dao.Client) dto.ClientListResponse {
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

func SupplierDAOtoDTO(supplier *dao.Supplier) dto.SupplierResponse {
	return dto.SupplierResponse{
		ID:          supplier.ID.String(),
		Name:        supplier.Name,
		PhoneNumber: supplier.PhoneNumber,
	}
}

func SupplierListDAOtoDTO(suppliers []*dao.Supplier) dto.SupplierListResponse {
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

func ProductDAOtoDTO(product *dao.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:             product.ID.String(),
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		LastUpdate:     product.LastUpdate.String(),
		SupplierID:     product.SupplierID.String(),
	}
}

func ProductListDAOtoDTO(products []*dao.Product) dto.ProductListResponse {
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
