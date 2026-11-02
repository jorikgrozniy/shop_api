package service

import (
	"context"
	"errors"
	"fmt"
	"shop_api/internal/entity"
	"shop_api/internal/repository"
	"shop_api/pkg/postgres"

	"github.com/google/uuid"
)

type SupplierService struct {
	supplierRepo   repository.SupplierRepository
	addressService *AddressService
	txManager      TxManager
}

func NewSupplierService(supplierRepo repository.SupplierRepository,
	addressService *AddressService, txManager TxManager) *SupplierService {
	return &SupplierService{
		supplierRepo:   supplierRepo,
		addressService: addressService,
		txManager:      txManager,
	}
}

func (s *SupplierService) AddSupplier(ctx context.Context, supplier *entity.Supplier, address *entity.Address) error {
	if supplier == nil {
		return fmt.Errorf("add supplier: %w", ErrNilEntity)
	}

	if len(supplier.Name) == 0 || len(supplier.Name) > 100 {
		return ErrInvalidNameLength
	}

	if len(supplier.PhoneNumber) == 0 || len(supplier.PhoneNumber) > 20 {
		return ErrInvalidPhoneLength
	}

	return s.txManager.WithTx(ctx, func(q postgres.Querier) error {
		addressID, err := s.addressService.MustGetAddressID(ctx, q, address)
		if err != nil {
			return fmt.Errorf("add supplier: %w", err)
		}
		supplier.AddressID = addressID

		if err := s.supplierRepo.Save(ctx, q, supplier); err != nil {
			return fmt.Errorf("%w: %w", ErrServerInternal, err)
		}

		return nil
	})
}

func (s *SupplierService) GetSupplier(ctx context.Context, id uuid.UUID) (*entity.Supplier, error) {
	supplier, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, fmt.Errorf("get supplier %s: %w", id, ErrSupplierNotFound)
		}
		return nil, fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return supplier, nil
}

func (s *SupplierService) RemoveSupplier(ctx context.Context, supplierID uuid.UUID) error {
	err := s.supplierRepo.RemoveByID(ctx, supplierID)

	switch {
	case err == nil:
		return nil
	case errors.Is(err, repository.ErrNoRows):
		return fmt.Errorf("remove supplier %s: %w", supplierID, ErrSupplierNotFound)
	case errors.Is(err, repository.ErrDependentEntity):
		return fmt.Errorf("remove supplier %s: %w", supplierID, ErrDependentEntity)
	default:
		return fmt.Errorf("%w: %w", ErrServerInternal, err)
	}
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
		return nil, 0, 0, fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return suppliers, limit, offset, nil
}

func (s *SupplierService) ChangeSupplierAddress(ctx context.Context, supplierID uuid.UUID, newAddress *entity.Address) error {
	return s.txManager.WithTx(ctx, func(q postgres.Querier) error {
		addressID, err := s.addressService.MustGetAddressID(ctx, q, newAddress)
		if err != nil {
			return fmt.Errorf("change address of supplier %s: %w", supplierID, err)
		}

		err = s.supplierRepo.UpdateAddress(ctx, q, supplierID, addressID)
		if err != nil {
			if errors.Is(err, repository.ErrNoRows) {
				return fmt.Errorf("change address of supplier %s: %w", supplierID, ErrSupplierNotFound)
			}
			return fmt.Errorf("%w: %w", ErrServerInternal, err)
		}

		return nil
	})
}
