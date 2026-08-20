package service

import (
	"context"
	"shop_api/internal/entity"
	"shop_api/internal/repository"

	"github.com/google/uuid"
)

type SupplierService struct {
	supplierRepo   repository.SupplierRepository
	addressService *AddressService
}

func NewSupplierService(supplierRepo repository.SupplierRepository,
	addressService *AddressService) *SupplierService {
	return &SupplierService{
		supplierRepo:   supplierRepo,
		addressService: addressService,
	}
}

func (s *SupplierService) AddSupplier(ctx context.Context, supplier *entity.Supplier, address *entity.Address) error {
	if len(supplier.Name) == 0 || len(supplier.Name) > 100 {
		return ErrInvalidNameLength
	}

	if len(supplier.PhoneNumber) == 0 || len(supplier.PhoneNumber) > 20 {
		return ErrInvalidPhoneLength
	}

	addressID, err := s.addressService.MustGetAddressID(ctx, address)
	if err != nil {
		return err
	}
	supplier.AddressID = addressID

	if err := s.supplierRepo.Save(ctx, supplier); err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *SupplierService) GetSupplier(ctx context.Context, id uuid.UUID) (*entity.Supplier, error) {
	supplier, err := s.supplierRepo.GetByID(ctx, id)

	if err == repository.ErrNoRows {
		return nil, ErrSupplierNotFound
	} else if err != nil {
		return nil, ErrServerInternal
	}

	return supplier, nil
}

func (s *SupplierService) RemoveSupplier(ctx context.Context, supplierID uuid.UUID) error {
	err := s.supplierRepo.RemoveByID(ctx, supplierID)

	switch err {
	case repository.ErrNoRows:
		return ErrSupplierNotFound
	case repository.ErrDependentEntity:
		return ErrDependentEntity
	}

	if err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *SupplierService) GetSuppliersWithParams(ctx context.Context, limit, offset int) ([]*entity.Supplier, int, int, error) {
	if limit < 1 || limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	suppliers, err := s.supplierRepo.GetWithParams(ctx, limit, offset)

	if err != nil {
		return nil, 0, 0, ErrServerInternal
	}

	return suppliers, limit, offset, nil
}

func (s *SupplierService) ChangeSupplierAddress(ctx context.Context, supplierID uuid.UUID, newAddress *entity.Address) error {
	addressID, err := s.addressService.MustGetAddressID(ctx, newAddress)

	if err != nil {
		return err
	}

	if err := s.supplierRepo.UpdateAddress(ctx, supplierID, addressID); err == repository.ErrNoRows {
		return ErrSupplierNotFound
	} else if err != nil {
		return ErrServerInternal
	}

	return nil
}
