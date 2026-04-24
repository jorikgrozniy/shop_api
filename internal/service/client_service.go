package service

import (
	"shop_api/internal/dao"
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

func (s *ClientService) AddClient(client *dao.Client, address *dao.Address) error {
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
	return s.clientRepo.RemoveByID(clientID)
}

func (s *ClientService) GetClientByName(name, surname string) (*dao.Client, error) {
	return s.clientRepo.GetByName(name, surname)
}

func (s *ClientService) GetAllClients(limit, offset int) ([]*dao.Client, error) {
	if limit < 1 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.clientRepo.GetAll(limit, offset)
}

func (s *ClientService) ChangeClientAddress(clientID uuid.UUID, newAddress *dao.Address) error {
	addressID, err := s.addressService.MustGetAddressID(newAddress)
	if err != nil {
		return err
	}
	return s.clientRepo.UpdateAddress(clientID, addressID)
}
