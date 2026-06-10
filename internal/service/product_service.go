package service

import (
	"shop_api/internal/entity"
	"shop_api/internal/repository"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepo     repository.ProductRepository
	supplierService *SupplierService
	imageService    *ImageService
}

func NewProductService(productRepo repository.ProductRepository,
	supplierService *SupplierService, imageService *ImageService) *ProductService {
	return &ProductService{
		productRepo:     productRepo,
		supplierService: supplierService,
		imageService:    imageService,
	}
}

func (s *ProductService) AddProduct(product *entity.Product, image *entity.Image) error {
	if len(product.Name) == 0 || len(product.Name) > 100 {
		return ErrInvalidNameLength
	}

	if len(product.Category) == 0 || len(product.Category) > 50 {
		return ErrInvalidCategoryLength
	}

	if product.Price < 0 {
		return ErrInvalidPrice
	}

	if product.AvailableStock <= 0 {
		return ErrInvalidAvailableStock
	}

	if _, err := s.supplierService.GetSupplier(*product.SupplierID); err != nil {
		return ErrSupplierNotFound
	}

	if image != nil {
		imageID, err := s.imageService.AddImage(image)
		if err != nil {
			return err
		}

		product.ImageID = &imageID
	}

	if err := s.productRepo.Save(product); err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ProductService) AddProductImage(productID uuid.UUID, image *entity.Image) error {
	imageID, err := s.imageService.AddImage(image)

	if err != nil {
		return err
	}

	if err := s.productRepo.UpdateImage(productID, imageID); err == repository.ErrNoRows {
		return ErrProductNotFound
	} else if err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ProductService) GetProductImage(productID uuid.UUID) (*entity.Image, error) {
	product, err := s.productRepo.GetByID(productID)

	if err == repository.ErrNoRows {
		return nil, ErrProductNotFound
	} else if err != nil {
		return nil, ErrServerInternal
	}

	if product.ImageID == nil {
		return nil, ErrImageNotFound
	}

	return s.imageService.GetImage(*product.ImageID)
}

func (s *ProductService) GetProduct(id uuid.UUID) (*entity.Product, error) {
	product, err := s.productRepo.GetByID(id)

	if err == repository.ErrNoRows {
		return nil, ErrProductNotFound
	} else if err != nil {
		return nil, ErrServerInternal
	}

	return product, nil
}

func (s *ProductService) DecreaseAvailableStock(productID uuid.UUID, value int) error {
	if value < 1 {
		return ErrInvalidAmount
	}

	if err := s.productRepo.UpdateAvailableStock(productID, -value); err == repository.ErrNoRows {
		return ErrProductNotFound
	} else if err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ProductService) RemoveProduct(productID uuid.UUID) error {
	err := s.productRepo.RemoveByID(productID)

	switch err {
	case repository.ErrNoRows:
		return ErrProductNotFound
	case repository.ErrDependentEntity:
		return ErrDependentEntity
	}

	if err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ProductService) GetProductsWithParams(limit, offset int) ([]*entity.Product, int, int, error) {
	if limit < 1 || limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	products, err := s.productRepo.GetWithParams(limit, offset)

	if err != nil {
		return nil, 0, 0, ErrServerInternal
	}

	return products, limit, offset, nil
}
