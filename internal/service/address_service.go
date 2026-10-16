package service

import (
	"context"
	"errors"
	"fmt"
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

func (s *AddressService) MustGetAddressID(ctx context.Context, address *entity.Address) (uuid.UUID, error) {
	if address == nil {
		return uuid.Nil, fmt.Errorf("resolve address: %w", ErrNilEntity)
	}

	id, found, err := s.findAddress(ctx, address)
	if err != nil {
		return uuid.Nil, err
	}

	if found {
		return id, nil
	}

	return s.addAddress(ctx, address)
}

func (s *AddressService) GetAddress(ctx context.Context, id uuid.UUID) (*entity.Address, error) {
	address, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, fmt.Errorf("get address %s: %w", id, ErrAddressNotFound)
		}
		return nil, fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return address, nil
}

func (s *AddressService) addAddress(ctx context.Context, address *entity.Address) (uuid.UUID, error) {
	if len(address.Country) == 0 || len(address.Country) > 100 {
		return uuid.Nil, ErrInvalidCountryLength
	}

	if len(address.City) == 0 || len(address.City) > 100 {
		return uuid.Nil, ErrInvalidCityLength
	}

	if len(address.Street) == 0 || len(address.Street) > 100 {
		return uuid.Nil, ErrInvalidStreetLength
	}

	id, err := s.addressRepo.Save(ctx, address)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return id, nil
}

func (s *AddressService) findAddress(ctx context.Context, address *entity.Address) (uuid.UUID, bool, error) {
	id, err := s.addressRepo.Find(ctx, address)
	if err == nil {
		return id, true, nil
	}

	if errors.Is(err, repository.ErrNoRows) {
		return uuid.Nil, false, nil
	}

	return uuid.Nil, false, fmt.Errorf("%w: %w", ErrServerInternal, err)
}
