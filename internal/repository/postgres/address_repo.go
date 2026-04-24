package postgres

import (
	"context"
	"shop_api/internal/dao"
	"shop_api/internal/repository"
	"shop_api/pkg/postgres"

	"github.com/google/uuid"
)

type addressRepoPostgres struct {
	pg *postgres.Postgres
}

func NewAddressRepoPostgres(pg *postgres.Postgres) repository.AddressRepository {
	return &addressRepoPostgres{
		pg: pg,
	}
}

func (r *addressRepoPostgres) Save(address *dao.Address) (uuid.UUID, error) {
	if address == nil {
		return uuid.Nil, ErrNilEntity
	}

	query := `
		INSERT INTO addresses (
			country, city, street
		) VALUES (
			$1, $2, $3
		)
		RETURNING id
	`

	var id uuid.UUID
	err := r.pg.Pool.QueryRow(context.Background(), query,
		address.Country, address.City, address.Street).Scan(&id)

	if err != nil {
		return uuid.Nil, ErrQueryExec
	}

	return id, nil
}

func (r *addressRepoPostgres) Find(address *dao.Address) (uuid.UUID, error) {
	query := `
		SELECT id
		FROM addresses
		WHERE country = $1 AND city = $2 AND street = $3
	`

	var id uuid.UUID
	err := r.pg.Pool.QueryRow(context.Background(), query,
		address.Country, address.City, address.Street).Scan(&id)

	if err != nil {
		return uuid.Nil, ErrQueryExec
	}

	return id, nil
}

func (r *addressRepoPostgres) GetByID(id uuid.UUID) (*dao.Address, error) {
	var address dao.Address
	query := `
		SELECT
			id, country, city, street
		FROM addresses
		WHERE id = $1
	`

	err := r.pg.Pool.QueryRow(context.Background(), query, id).Scan(
		&address.ID, &address.Country, &address.City, &address.Street,
	)

	if err != nil {
		return nil, ErrQueryExec
	}

	return &address, nil
}
