package postgres

import (
	"context"
	"shop_api/internal/dao"
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

func (r *supplierRepoPostgres) Save(supplier *dao.Supplier) error {
	if supplier == nil {
		return ErrNilEntity
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
		return ErrQueryExec
	}

	return nil
}

func (r *supplierRepoPostgres) RemoveByID(id uuid.UUID) error {
	query := `
		DELETE FROM suppliers
		WHERE id = $1
	`

	_, err := r.pg.Pool.Exec(context.Background(), query, id)

	if err != nil {
		return ErrQueryExec
	}

	return nil
}

func (r *supplierRepoPostgres) GetByID(id uuid.UUID) (*dao.Supplier, error) {
	var supplier dao.Supplier
	query := `
		SELECT
			id, name, address_id, phone_number
		FROM suppliers
		WHERE id = $1
	`

	err := r.pg.Pool.QueryRow(context.Background(), query, id).Scan(
		&supplier.ID, &supplier.Name, &supplier.AddressID, &supplier.PhoneNumber,
	)

	if err != nil {
		return nil, ErrQueryExec
	}

	return &supplier, nil
}

func (r *supplierRepoPostgres) GetAll(limit, offset int) ([]*dao.Supplier, error) {
	query := `
		SELECT
			id, name, address_id, phone_number
		FROM suppliers
		LIMIT $1
		OFFSET $2
	`

	rows, err := r.pg.Pool.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, ErrQueryExec
	}
	defer rows.Close()

	var suppliers []*dao.Supplier
	for rows.Next() {
		var supplier dao.Supplier
		err := rows.Scan(
			&supplier.ID, &supplier.Name,
			&supplier.AddressID, &supplier.PhoneNumber,
		)
		if err != nil {
			return nil, ErrRowScan
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

	_, err := r.pg.Pool.Exec(context.Background(), query, supplierID, addressID)

	if err != nil {
		return ErrQueryExec
	}

	return nil
}
