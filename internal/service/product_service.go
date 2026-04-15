package service

import (
	"shop_api/internal/dao"
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

func (s *ProductService) AddProduct(product *dao.Product, image *dao.Image) error {
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

	if _, err := s.supplierService.GetSupplier(product.SupplierID); err != nil {
		return ErrSupplierNotFound
	}

	if image != nil {
		imageID, err := s.imageService.AddImage(image.Image)
		if err != nil {
			return err
		}

		product.ImageID = imageID
	}

	if err := s.productRepo.Save(product); err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ProductService) AddProductImage(productID uuid.UUID, image []byte) error {

}

func (s *ProductService) GetProductImage(productID uuid.UUID) ([]byte, error) {

}

func (s *ProductService) GetProduct(id uuid.UUID) (*dao.Product, error) {
	return s.productRepo.GetByID(id)
}

func (s *ProductService) DecreaseAvailableStock(productID uuid.UUID, value int) error {
	return s.productRepo.UpdateAvailableStock(productID, -value)
}

func (s *ProductService) RemoveProduct(productID uuid.UUID) error {
	return s.productRepo.RemoveByID(productID)
}

func (s *ProductService) GetAllProducts(limit, offset int) ([]*dao.Product, error) {
	return s.productRepo.GetAll(limit, offset)
}
