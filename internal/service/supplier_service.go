package service

import (
	"shop_api/internal/dao"
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

func (s *SupplierService) AddSupplier(supplier *dao.Supplier, address *dao.Address) error {
	if len(supplier.Name) == 0 || len(supplier.Name) > 100 {
		return ErrInvalidNameLength
	}

	if len(supplier.PhoneNumber) == 0 || len(supplier.PhoneNumber) > 20 {
		return ErrInvalidPhoneLength
	}

	addressID, err := s.addressService.MustGetAddressID(address)
	if err != nil {
		return err
	}
	supplier.AddressID = addressID

	if err := s.supplierRepo.Save(supplier); err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *SupplierService) GetSupplier(id uuid.UUID) (*dao.Supplier, error) {
	return s.supplierRepo.GetByID(id)
}

func (s *SupplierService) RemoveSupplier(supplierID uuid.UUID) error {
	return s.supplierRepo.RemoveByID(supplierID)
}

func (s *SupplierService) GetAllSuppliers(limit, offset int) ([]*dao.Supplier, error) {
	return s.supplierRepo.GetAll(limit, offset)
}

func (s *SupplierService) ChangeSupplierAddress(supplierID uuid.UUID, newAddress *dao.Address) error {
	return s.supplierRepo.UpdateAddress(supplierID, newAddress)
}
