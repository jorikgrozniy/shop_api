package service

import (
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

func (s *ClientService) AddClient(client *entity.Client, address *entity.Address) error {
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

	addressID, err := s.addressService.MustGetAddressID(address)
	if err != nil {
		return err
	}
	client.AddressID = addressID

	if err := s.clientRepo.Save(client); err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ClientService) RemoveClient(clientID uuid.UUID) error {
	err := s.clientRepo.RemoveByID(clientID)

	switch err {
	case repository.ErrNoRows:
		return ErrClientNotFound
	case repository.ErrDependentEntity:
		return ErrDependentEntity
	}

	if err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ClientService) GetClientsWithParams(name, surname *string, limit, offset int) ([]*entity.Client, int, int, error) {
	if limit < 1 || limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	if *name == "" {
		name = nil
	}

	if *surname == "" {
		surname = nil
	}

	clients, err := s.clientRepo.GetWithParams(name, surname, limit, offset)
	if err != nil {
		return nil, 0, 0, ErrServerInternal
	}

	return clients, limit, offset, nil
}

func (s *ClientService) ChangeClientAddress(clientID uuid.UUID, newAddress *entity.Address) error {
	addressID, err := s.addressService.MustGetAddressID(newAddress)

	if err != nil {
		return err
	}

	if err := s.clientRepo.UpdateAddress(clientID, addressID); err == repository.ErrNoRows {
		return ErrClientNotFound
	} else if err != nil {
		return ErrServerInternal
	}

	return nil
}
