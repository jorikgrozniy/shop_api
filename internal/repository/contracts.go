package repository

import (
	"shop_api/internal/entity"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Save(client *entity.Client) error
	RemoveByID(id uuid.UUID) error
	GetWithParams(name, surname *string, limit, offset int) ([]*entity.Client, error)
	UpdateAddress(clientID, addressID uuid.UUID) error
}

type AddressRepository interface {
	Save(address *entity.Address) (uuid.UUID, error)
	Find(address *entity.Address) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*entity.Address, error)
}

type ImageRepository interface {
	Save(image *entity.Image) (uuid.UUID, error)
	RemoveByID(id uuid.UUID) error
	GetByID(id uuid.UUID) (*entity.Image, error)
	Update(id uuid.UUID, newImage []byte) error
}

type SupplierRepository interface {
	Save(supplier *entity.Supplier) error
	RemoveByID(id uuid.UUID) error
	GetByID(id uuid.UUID) (*entity.Supplier, error)
	GetWithParams(limit, offset int) ([]*entity.Supplier, error)
	UpdateAddress(supplierID, addressID uuid.UUID) error
}

type ProductRepository interface {
	Save(product *entity.Product) error
	UpdateAvailableStock(id uuid.UUID, value int) error
	UpdateImage(productID, imageID uuid.UUID) error
	GetByID(id uuid.UUID) (*entity.Product, error)
	GetWithParams(limit, offset int) ([]*entity.Product, error)
	RemoveByID(id uuid.UUID) error
}
