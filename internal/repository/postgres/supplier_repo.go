package postgres

import (
	"context"
	"shop_api/internal/entity"
	"shop_api/internal/repository"
	"shop_api/pkg/postgres"

	"github.com/google/uuid"
)

type supplierRepoPostgres struct {
	pg *postgres.Postgres
}

func NewSupplierRepoPostgres(pg *postgres.Postgres) repository.SupplierRepository {
	return &supplierRepoPostgres{
		pg: pg,
	}
}

func (r *supplierRepoPostgres) Save(supplier *entity.Supplier) error {
	if supplier == nil {
		return repository.ErrNilEntity
	}

	query := `
		INSERT INTO suppliers (
			name, address_id, phone_number
		) VALUES (
			$1, $2, $3
		)
	`

	_, err := r.pg.Pool.Exec(context.Background(), query,
		supplier.Name, supplier.AddressID, supplier.PhoneNumber,
	)

	if err != nil {
		return repository.ErrQueryExec
	}

	return nil
}

func (r *supplierRepoPostgres) RemoveByID(id uuid.UUID) error {
	query := `
		DELETE FROM suppliers
		WHERE id = $1
	`

	cmdTag, err := r.pg.Pool.Exec(context.Background(), query, id)

	if isFKViolation(err) {
		return repository.ErrDependentEntity
	} else if err != nil {
		return repository.ErrQueryExec
	}

	if cmdTag.RowsAffected() == 0 {
		return repository.ErrNoRows
	}

	return nil
}

func (r *supplierRepoPostgres) GetByID(id uuid.UUID) (*entity.Supplier, error) {
	var supplier entity.Supplier
	query := `
		SELECT
			id, name, address_id, phone_number
		FROM suppliers
		WHERE id = $1
	`

	err := r.pg.Pool.QueryRow(context.Background(), query, id).Scan(
		&supplier.ID, &supplier.Name, &supplier.AddressID, &supplier.PhoneNumber,
	)

	if err == postgres.ErrNoRows {
		return nil, repository.ErrNoRows
	} else if err != nil {
		return nil, repository.ErrQueryExec
	}

	return &supplier, nil
}

func (r *supplierRepoPostgres) GetWithParams(limit, offset int) ([]*entity.Supplier, error) {
	query := `
		SELECT
			id, name, address_id, phone_number
		FROM suppliers
		LIMIT $1
		OFFSET $2
	`

	rows, err := r.pg.Pool.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, repository.ErrQueryExec
	}
	defer rows.Close()

	var suppliers []*entity.Supplier
	for rows.Next() {
		var supplier entity.Supplier
		err := rows.Scan(
			&supplier.ID, &supplier.Name,
			&supplier.AddressID, &supplier.PhoneNumber,
		)
		if err != nil {
			return nil, repository.ErrRowScan
		}
		suppliers = append(suppliers, &supplier)
	}

	return suppliers, nil
}

func (r *supplierRepoPostgres) UpdateAddress(supplierID, addressID uuid.UUID) error {
	query := `
		UPDATE suppliers
		SET address_id = $2
		WHERE id = $1
	`

	cmdTag, err := r.pg.Pool.Exec(context.Background(), query, supplierID, addressID)

	if err != nil {
		return repository.ErrQueryExec
	}

	if cmdTag.RowsAffected() == 0 {
		return repository.ErrNoRows
	}

	return nil
}
