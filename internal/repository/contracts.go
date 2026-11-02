package repository

import (
	"context"
	"shop_api/internal/entity"
	"shop_api/pkg/postgres"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Save(ctx context.Context, q postgres.Querier, client *entity.Client) error
	RemoveByID(ctx context.Context, id uuid.UUID) error
	GetWithParams(ctx context.Context, name, surname *string, limit, offset int) ([]*entity.Client, error)
	UpdateAddress(ctx context.Context, q postgres.Querier, clientID, addressID uuid.UUID) error
}

type AddressRepository interface {
	Save(ctx context.Context, q postgres.Querier, address *entity.Address) (uuid.UUID, error)
	Find(ctx context.Context, q postgres.Querier, address *entity.Address) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Address, error)
}

type ImageRepository interface {
	Save(ctx context.Context, q postgres.Querier, image *entity.Image) (uuid.UUID, error)
	RemoveByID(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Image, error)
	Update(ctx context.Context, id uuid.UUID, newImage []byte) error
}

type SupplierRepository interface {
	Save(ctx context.Context, q postgres.Querier, supplier *entity.Supplier) error
	RemoveByID(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Supplier, error)
	GetWithParams(ctx context.Context, limit, offset int) ([]*entity.Supplier, error)
	UpdateAddress(ctx context.Context, q postgres.Querier, supplierID, addressID uuid.UUID) error
}

type ProductRepository interface {
	Save(ctx context.Context, q postgres.Querier, product *entity.Product) error
	UpdateAvailableStock(ctx context.Context, id uuid.UUID, value int) error
	UpdateImage(ctx context.Context, q postgres.Querier, productID, imageID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	GetWithParams(ctx context.Context, limit, offset int) ([]*entity.Product, error)
	RemoveByID(ctx context.Context, id uuid.UUID) error
}
