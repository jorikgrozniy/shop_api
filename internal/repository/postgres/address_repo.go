package postgres

import (
	"context"
	"shop_api/internal/entity"
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

func (r *addressRepoPostgres) Save(ctx context.Context, address *entity.Address) (uuid.UUID, error) {
	if address == nil {
		return uuid.Nil, repository.ErrNilEntity
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
	err := r.pg.Pool.QueryRow(ctx, query,
		address.Country, address.City, address.Street).Scan(&id)

	if err == postgres.ErrNoRows {
		return uuid.Nil, repository.ErrNoRows
	} else if err != nil {
		return uuid.Nil, repository.ErrQueryExec
	}

	return id, nil
}

func (r *addressRepoPostgres) Find(ctx context.Context, address *entity.Address) (uuid.UUID, error) {
	query := `
		SELECT id
		FROM addresses
		WHERE country = $1 AND city = $2 AND street = $3
	`

	var id uuid.UUID
	err := r.pg.Pool.QueryRow(ctx, query,
		address.Country, address.City, address.Street).Scan(&id)

	if err == postgres.ErrNoRows {
		return uuid.Nil, repository.ErrNoRows
	} else if err != nil {
		return uuid.Nil, repository.ErrQueryExec
	}

	return id, nil
}

func (r *addressRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.Address, error) {
	var address entity.Address
	query := `
		SELECT
			id, country, city, street
		FROM addresses
		WHERE id = $1
	`

	err := r.pg.Pool.QueryRow(ctx, query, id).Scan(
		&address.ID, &address.Country, &address.City, &address.Street,
	)

	if err == postgres.ErrNoRows {
		return nil, repository.ErrNoRows
	} else if err != nil {
		return nil, repository.ErrQueryExec
	}

	return &address, nil
}
