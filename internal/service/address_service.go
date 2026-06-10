package service

import (
	"shop_api/internal/entity"
	"shop_api/internal/repository"

	"github.com/google/uuid"
)

type AddressService struct {
	addressRepo repository.AddressRepository
}

func NewAddressService(addressRepo repository.AddressRepository) *AddressService {
	return &AddressService{
		addressRepo: addressRepo,
	}
}

func (s *AddressService) MustGetAddressID(address *entity.Address) (uuid.UUID, error) {
	if id, found := s.findAddress(address); found {
		return id, nil
	}

	if id, err := s.addAddress(address); err != nil {
		return uuid.Nil, err
	} else {
		return id, nil
	}
}

func (s *AddressService) GetAddress(id uuid.UUID) (*entity.Address, error) {
	address, err := s.addressRepo.GetByID(id)

	if err == repository.ErrNoRows {
		return nil, ErrAddressNotFound
	} else if err != nil {
		return nil, ErrServerInternal
	}

	return address, nil
}

func (s *AddressService) addAddress(address *entity.Address) (uuid.UUID, error) {
	if len(address.Country) == 0 || len(address.Country) > 100 {
		return uuid.Nil, ErrInvalidCountryLength
	}

	if len(address.City) == 0 || len(address.City) > 100 {
		return uuid.Nil, ErrInvalidCityLength
	}

	if len(address.Street) == 0 || len(address.Street) > 100 {
		return uuid.Nil, ErrInvalidStreetLength
	}

	id, err := s.addressRepo.Save(address)
	if err != nil {
		return uuid.Nil, ErrServerInternal
	}

	return id, nil
}

func (s *AddressService) findAddress(address *entity.Address) (uuid.UUID, bool) {
	if id, err := s.addressRepo.Find(address); err != nil {
		return uuid.Nil, false
	} else {
		return id, true
	}
}
