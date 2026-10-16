package service

import (
	"context"
	"errors"
	"fmt"
	"shop_api/internal/entity"
	"shop_api/internal/repository"
	"time"

	"github.com/google/uuid"
)

type ClientService struct {
	clientRepo     repository.ClientRepository
	addressService *AddressService
}

func NewClientService(clientRepo repository.ClientRepository,
	addressService *AddressService) *ClientService {
	return &ClientService{
		clientRepo:     clientRepo,
		addressService: addressService,
	}
}

func (s *ClientService) AddClient(ctx context.Context, client *entity.Client, address *entity.Address) error {
	if client == nil {
		return fmt.Errorf("add client: %w", ErrNilEntity)
	}

	if len(client.Name) == 0 || len(client.Name) > 100 {
		return ErrInvalidNameLength
	}

	if len(client.Surname) == 0 || len(client.Surname) > 100 {
		return ErrInvalidSurnameLength
	}

	now := time.Now()
	if client.Birthdate.AddDate(18, 0, 0).After(now) {
		return ErrNotAdult
	}

	if client.Gender != "male" && client.Gender != "female" {
		return ErrInvalidGender
	}

	addressID, err := s.addressService.MustGetAddressID(ctx, address)
	if err != nil {
		return fmt.Errorf("add client: %w", err)
	}
	client.AddressID = addressID

	if err := s.clientRepo.Save(ctx, client); err != nil {
		return fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return nil
}

func (s *ClientService) RemoveClient(ctx context.Context, clientID uuid.UUID) error {
	err := s.clientRepo.RemoveByID(ctx, clientID)

	switch {
	case err == nil:
		return nil
	case errors.Is(err, repository.ErrNoRows):
		return fmt.Errorf("remove client %s: %w", clientID, ErrClientNotFound)
	case errors.Is(err, repository.ErrDependentEntity):
		return fmt.Errorf("remove client %s: %w", clientID, ErrDependentEntity)
	default:
		return fmt.Errorf("%w: %w", ErrServerInternal, err)
	}
}

func (s *ClientService) GetClientsWithParams(
	ctx context.Context, name, surname *string, limit, offset int) ([]*entity.Client, int, int, error) {
	if limit < 1 || limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	if name == nil || *name == "" {
		name = nil
	}

	if surname == nil || *surname == "" {
		surname = nil
	}

	clients, err := s.clientRepo.GetWithParams(ctx, name, surname, limit, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return clients, limit, offset, nil
}

func (s *ClientService) ChangeClientAddress(ctx context.Context, clientID uuid.UUID, newAddress *entity.Address) error {
	addressID, err := s.addressService.MustGetAddressID(ctx, newAddress)
	if err != nil {
		return fmt.Errorf("change address of client %s: %w", clientID, err)
	}

	err = s.clientRepo.UpdateAddress(ctx, clientID, addressID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return fmt.Errorf("change address of client %s: %w", clientID, ErrClientNotFound)
		}
		return fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return nil
}
