package service

import (
	"context"
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

func (s *ProductService) AddProduct(ctx context.Context, product *entity.Product, image *entity.Image) error {
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

	if _, err := s.supplierService.GetSupplier(ctx, *product.SupplierID); err != nil {
		return ErrSupplierNotFound
	}

	if image != nil {
		imageID, err := s.imageService.AddImage(ctx, image)
		if err != nil {
			return err
		}

		product.ImageID = &imageID
	}

	if err := s.productRepo.Save(ctx, product); err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ProductService) AddProductImage(ctx context.Context, productID uuid.UUID, image *entity.Image) error {
	imageID, err := s.imageService.AddImage(ctx, image)

	if err != nil {
		return err
	}

	if err := s.productRepo.UpdateImage(ctx, productID, imageID); err == repository.ErrNoRows {
		return ErrProductNotFound
	} else if err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ProductService) GetProductImage(ctx context.Context, productID uuid.UUID) (*entity.Image, error) {
	product, err := s.productRepo.GetByID(ctx, productID)

	if err == repository.ErrNoRows {
		return nil, ErrProductNotFound
	} else if err != nil {
		return nil, ErrServerInternal
	}

	if product.ImageID == nil {
		return nil, ErrImageNotFound
	}

	return s.imageService.GetImage(ctx, *product.ImageID)
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)

	if err == repository.ErrNoRows {
		return nil, ErrProductNotFound
	} else if err != nil {
		return nil, ErrServerInternal
	}

	return product, nil
}

func (s *ProductService) DecreaseAvailableStock(ctx context.Context, productID uuid.UUID, value int) error {
	if value < 1 {
		return ErrInvalidAmount
	}

	if err := s.productRepo.UpdateAvailableStock(ctx, productID, -value); err == repository.ErrNoRows {
		return ErrProductNotFound
	} else if err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ProductService) RemoveProduct(ctx context.Context, productID uuid.UUID) error {
	err := s.productRepo.RemoveByID(ctx, productID)

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

func (s *ProductService) GetProductsWithParams(ctx context.Context, limit, offset int) ([]*entity.Product, int, int, error) {
	if limit < 1 || limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	products, err := s.productRepo.GetWithParams(ctx, limit, offset)

	if err != nil {
		return nil, 0, 0, ErrServerInternal
	}

	return products, limit, offset, nil
}
