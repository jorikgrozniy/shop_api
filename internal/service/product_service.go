package service

import (
	"context"
	"errors"
	"fmt"
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
	if product == nil {
		return fmt.Errorf("add product: %w", ErrNilEntity)
	}

	if product.SupplierID == nil {
		return ErrInvalidSupplierID
	}

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
		return fmt.Errorf("add product: %w", err)
	}

	if image != nil {
		imageID, err := s.imageService.AddImage(ctx, image)
		if err != nil {
			return fmt.Errorf("add product: %w", err)
		}

		product.ImageID = &imageID
	}

	if err := s.productRepo.Save(ctx, product); err != nil {
		return fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return nil
}

func (s *ProductService) AddProductImage(ctx context.Context, productID uuid.UUID, image *entity.Image) error {
	imageID, err := s.imageService.AddImage(ctx, image)
	if err != nil {
		return fmt.Errorf("add image to product %s: %w", productID, err)
	}

	err = s.productRepo.UpdateImage(ctx, productID, imageID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return fmt.Errorf("add image to product %s: %w", productID, ErrProductNotFound)
		}
		return fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return nil
}

func (s *ProductService) GetProductImage(ctx context.Context, productID uuid.UUID) (*entity.Image, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, fmt.Errorf("get image of product %s: %w", productID, ErrProductNotFound)
		}
		return nil, fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	if product.ImageID == nil {
		return nil, fmt.Errorf("get image of product %s: %w", productID, ErrImageNotFound)
	}

	image, err := s.imageService.GetImage(ctx, *product.ImageID)
	if err != nil {
		return nil, fmt.Errorf("get image of product %s: %w", productID, err)
	}

	return image, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, fmt.Errorf("get product %s: %w", id, ErrProductNotFound)
		}
		return nil, fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return product, nil
}

func (s *ProductService) DecreaseAvailableStock(ctx context.Context, productID uuid.UUID, value int) error {
	if value < 1 {
		return ErrInvalidAmount
	}

	err := s.productRepo.UpdateAvailableStock(ctx, productID, -value)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return fmt.Errorf("decrease stock of product %s: %w", productID, ErrProductNotFound)
		}
		return fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return nil
}

func (s *ProductService) RemoveProduct(ctx context.Context, productID uuid.UUID) error {
	err := s.productRepo.RemoveByID(ctx, productID)

	switch {
	case err == nil:
		return nil
	case errors.Is(err, repository.ErrNoRows):
		return fmt.Errorf("remove product %s: %w", productID, ErrProductNotFound)
	case errors.Is(err, repository.ErrDependentEntity):
		return fmt.Errorf("remove product %s: %w", productID, ErrDependentEntity)
	default:
		return fmt.Errorf("%w: %w", ErrServerInternal, err)
	}
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
		return nil, 0, 0, fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return products, limit, offset, nil
}
