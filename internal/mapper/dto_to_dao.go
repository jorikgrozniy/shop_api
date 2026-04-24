package mapper

import (
	"shop_api/internal/dao"
	"shop_api/internal/dto"
	"time"

	"github.com/google/uuid"
)

func AddressDTOtoDAO(req dto.AddressRequest) *dao.Address {
	return &dao.Address{
		Country: req.Country,
		City:    req.City,
		Street:  req.Street,
	}
}

func ImageDTOtoDAO(req dto.ImageRequest) *dao.Image {
	return &dao.Image{
		Image: []byte(req.Image),
	}
}

func ClientDTOtoDAO(req dto.ClientRequest) (*dao.Client, error) {
	t, err := time.Parse("2006-01-02", req.Birthdate)
	if err != nil {
		return nil, err
	}

	return &dao.Client{
		Name:      req.Name,
		Surname:   req.Surname,
		Birthdate: t,
		Gender:    req.Gender,
	}, nil
}

func SupplierDTOtoDAO(req dto.SupplierRequest) *dao.Supplier {
	return &dao.Supplier{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}
}

func ProductDTOtoDAO(req dto.ProductRequest) (*dao.Product, error) {
	id, err := uuid.Parse(req.SupplierID)
	if err != nil {
		return nil, err
	}

	return &dao.Product{
		Name:           req.Name,
		Category:       req.Category,
		Price:          req.Price,
		AvailableStock: req.AvailableStock,
		SupplierID:     id,
	}, nil
}
