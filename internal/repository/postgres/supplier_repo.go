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

type supplierRepoPostgres struct {
	pg *postgres.Postgres
}

func NewSupplierRepoPostgres(pg *postgres.Postgres) repository.SupplierRepository {
	return &supplierRepoPostgres{
		pg: pg,
	}
}

func (r *supplierRepoPostgres) Save(ctx context.Context, supplier *entity.Supplier) error {
	if supplier == nil {
		return fmt.Errorf("save supplier: %w", repository.ErrNilEntity)
	}

	query := `
		INSERT INTO suppliers (
			name, address_id, phone_number
		) VALUES (
			$1, $2, $3
		)
	`

	_, err := r.pg.Pool.Exec(ctx, query,
		supplier.Name, supplier.AddressID, supplier.PhoneNumber,
	)

	if err != nil {
		return fmt.Errorf("save supplier: %w: %w", repository.ErrQueryExec, err)
	}

	return nil
}

func (r *supplierRepoPostgres) RemoveByID(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM suppliers
		WHERE id = $1
	`

	cmdTag, err := r.pg.Pool.Exec(ctx, query, id)

	if isFKViolation(err) {
		return fmt.Errorf("remove supplier %s: %w: %w", id, repository.ErrDependentEntity, err)
	} else if err != nil {
		return fmt.Errorf("remove supplier %s: %w: %w", id, repository.ErrQueryExec, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("remove supplier %s: %w", id, repository.ErrNoRows)
	}

	return nil
}

func (r *supplierRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.Supplier, error) {
	var supplier entity.Supplier
	query := `
		SELECT
			id, name, address_id, phone_number
		FROM suppliers
		WHERE id = $1
	`

	err := r.pg.Pool.QueryRow(ctx, query, id).Scan(
		&supplier.ID, &supplier.Name, &supplier.AddressID, &supplier.PhoneNumber,
	)

	if errors.Is(err, postgres.ErrNoRows) {
		return nil, fmt.Errorf("get supplier %s: %w", id, repository.ErrNoRows)
	} else if err != nil {
		return nil, fmt.Errorf("get supplier %s: %w: %w", id, repository.ErrQueryExec, err)
	}

	return &supplier, nil
}

func (r *supplierRepoPostgres) GetWithParams(ctx context.Context, limit, offset int) ([]*entity.Supplier, error) {
	query := `
		SELECT
			id, name, address_id, phone_number
		FROM suppliers
		LIMIT $1
		OFFSET $2
	`

	rows, err := r.pg.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list suppliers: %w: %w", repository.ErrQueryExec, err)
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
			return nil, fmt.Errorf("list suppliers: %w: %w", repository.ErrRowScan, err)
		}
		suppliers = append(suppliers, &supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list suppliers: %w: %w", repository.ErrQueryExec, err)
	}

	return suppliers, nil
}

func (r *supplierRepoPostgres) UpdateAddress(ctx context.Context, supplierID, addressID uuid.UUID) error {
	query := `
		UPDATE suppliers
		SET address_id = $2
		WHERE id = $1
	`

	cmdTag, err := r.pg.Pool.Exec(ctx, query, supplierID, addressID)

	if err != nil {
		return fmt.Errorf("update address of supplier %s: %w: %w", supplierID, repository.ErrQueryExec, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("update address of supplier %s: %w", supplierID, repository.ErrNoRows)
	}

	return nil
}
