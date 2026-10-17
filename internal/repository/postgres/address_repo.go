package postgres

import (
	"context"
	"errors"
	"fmt"
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

func (r *addressRepoPostgres) Save(ctx context.Context, q postgres.Querier, address *entity.Address) (uuid.UUID, error) {
	if address == nil {
		return uuid.Nil, fmt.Errorf("save address: %w", repository.ErrNilEntity)
	}

	query := `
		INSERT INTO addresses (
			country, city, street
		) VALUES (
			$1, $2, $3
		)
		ON CONFLICT (country, city, street)
		DO UPDATE SET country = addresses.country
		RETURNING id
	`

	var id uuid.UUID
	err := q.QueryRow(ctx, query,
		address.Country, address.City, address.Street).Scan(&id)

	if errors.Is(err, postgres.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("save address: %w", repository.ErrNoRows)
	} else if err != nil {
		return uuid.Nil, fmt.Errorf("save address: %w: %w", repository.ErrQueryExec, err)
	}

	return id, nil
}

func (r *addressRepoPostgres) Find(ctx context.Context, q postgres.Querier, address *entity.Address) (uuid.UUID, error) {
	if address == nil {
		return uuid.Nil, fmt.Errorf("find address: %w", repository.ErrNilEntity)
	}

	query := `
		SELECT id
		FROM addresses
		WHERE country = $1 AND city = $2 AND street = $3
	`

	var id uuid.UUID
	err := q.QueryRow(ctx, query,
		address.Country, address.City, address.Street).Scan(&id)

	if errors.Is(err, postgres.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("find address: %w", repository.ErrNoRows)
	} else if err != nil {
		return uuid.Nil, fmt.Errorf("find address: %w: %w", repository.ErrQueryExec, err)
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

	if errors.Is(err, postgres.ErrNoRows) {
		return nil, fmt.Errorf("get address %s: %w", id, repository.ErrNoRows)
	} else if err != nil {
		return nil, fmt.Errorf("get address %s: %w: %w", id, repository.ErrQueryExec, err)
	}

	return &address, nil
}
