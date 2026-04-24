package repository

import (
	"shop_api/internal/dao"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Save(client *dao.Client) error
	RemoveByID(id uuid.UUID) error
	GetByID(id uuid.UUID) (*dao.Client, error)
	GetByName(name, surname string) (*dao.Client, error)
	GetAll(limit, offset int) ([]*dao.Client, error)
	UpdateAddress(clientID, addressID uuid.UUID) error
}

type AddressRepository interface {
	Save(address *dao.Address) (uuid.UUID, error)
	Find(address *dao.Address) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*dao.Address, error)
}

type ImageRepository interface {
	Save(image []byte) (uuid.UUID, error)
	RemoveByID(id uuid.UUID) error
	GetByID(id uuid.UUID) ([]byte, error)
	Update(id uuid.UUID, newImage []byte) error
}

type SupplierRepository interface {
	Save(supplier *dao.Supplier) error
	RemoveByID(id uuid.UUID) error
	GetByID(id uuid.UUID) (*dao.Supplier, error)
	GetAll(limit, offset int) ([]*dao.Supplier, error)
	UpdateAddress(supplierID, addressID uuid.UUID) error
}

type ProductRepository interface {
	Save(product *dao.Product) error
	UpdateAvailableStock(id uuid.UUID, value int) error
	UpdateImage(productID, imageID uuid.UUID) error
	GetByID(id uuid.UUID) (*dao.Product, error)
	GetAll(limit, offset int) ([]*dao.Product, error)
	RemoveByID(id uuid.UUID) error
}
