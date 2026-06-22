package repository

import (
	"context"
	"shop_api/internal/entity"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Save(ctx context.Context, client *entity.Client) error
	RemoveByID(ctx context.Context, id uuid.UUID) error
	GetWithParams(ctx context.Context, name, surname *string, limit, offset int) ([]*entity.Client, error)
	UpdateAddress(ctx context.Context, clientID, addressID uuid.UUID) error
}

type AddressRepository interface {
	Save(ctx context.Context, address *entity.Address) (uuid.UUID, error)
	Find(ctx context.Context, address *entity.Address) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Address, error)
}

type ImageRepository interface {
	Save(ctx context.Context, image *entity.Image) (uuid.UUID, error)
	RemoveByID(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Image, error)
	Update(ctx context.Context, id uuid.UUID, newImage []byte) error
}

type SupplierRepository interface {
	Save(ctx context.Context, supplier *entity.Supplier) error
	RemoveByID(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Supplier, error)
	GetWithParams(ctx context.Context, limit, offset int) ([]*entity.Supplier, error)
	UpdateAddress(ctx context.Context, supplierID, addressID uuid.UUID) error
}

type ProductRepository interface {
	Save(ctx context.Context, product *entity.Product) error
	UpdateAvailableStock(ctx context.Context, id uuid.UUID, value int) error
	UpdateImage(ctx context.Context, productID, imageID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	GetWithParams(ctx context.Context, limit, offset int) ([]*entity.Product, error)
	RemoveByID(ctx context.Context, id uuid.UUID) error
}
