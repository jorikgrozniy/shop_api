package mapper

import (
	"shop_api/internal/controller/restapi/v1/dto"
	"shop_api/internal/entity"
	"time"

	"github.com/google/uuid"
)

func AddressDTOtoDAO(req dto.AddressRequest) *entity.Address {
	return &entity.Address{
		Country: req.Country,
		City:    req.City,
		Street:  req.Street,
	}
}

func ImageDTOtoDAO(req dto.ImageRequest) *entity.Image {
	return &entity.Image{
		Image: []byte(req.JPGbase64),
	}
}

func ClientDTOtoDAO(req dto.ClientRequest) (*entity.Client, error) {
	t, err := time.Parse("2006-01-02", req.Birthdate)
	if err != nil {
		return nil, err
	}

	return &entity.Client{
		Name:      req.Name,
		Surname:   req.Surname,
		Birthdate: t,
		Gender:    req.Gender,
	}, nil
}

func SupplierDTOtoDAO(req dto.SupplierRequest) *entity.Supplier {
	return &entity.Supplier{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}
}

func ProductDTOtoDAO(req dto.ProductRequest) (*entity.Product, error) {
	id, err := uuid.Parse(req.SupplierID)
	if err != nil {
		return nil, err
	}

	return &entity.Product{
		Name:           req.Name,
		Category:       req.Category,
		Price:          req.Price,
		AvailableStock: req.AvailableStock,
		SupplierID:     &id,
	}, nil
}
