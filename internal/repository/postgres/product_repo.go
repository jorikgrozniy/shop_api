package postgres

import (
	"context"
	"shop_api/internal/entity"
	"shop_api/internal/repository"
	"shop_api/pkg/postgres"

	"github.com/google/uuid"
)

type productRepoPostgres struct {
	pg *postgres.Postgres
}

func NewProductRepoPostgres(pg *postgres.Postgres) repository.ProductRepository {
	return &productRepoPostgres{
		pg: pg,
	}
}

func (r *productRepoPostgres) Save(product *entity.Product) error {
	if product == nil {
		return repository.ErrNilEntity
	}

	query := `
		INSERT INTO products (
			name, category, price, available_stock,
			last_update_date, supplier_id, image_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	_, err := r.pg.Pool.Exec(context.Background(), query,
		product.Name, product.Category, product.Price, product.AvailableStock,
		product.LastUpdate, product.SupplierID, product.ImageID,
	)

	if err != nil {
		return repository.ErrQueryExec
	}

	return nil
}

func (r *productRepoPostgres) UpdateAvailableStock(id uuid.UUID, value int) error {
	query := `
		UPDATE products
		SET available_stock = available_stock + $2
		WHERE id = $1
	`

	cmdTag, err := r.pg.Pool.Exec(context.Background(), query, id, value)

	if err != nil {
		return repository.ErrQueryExec
	}

	if cmdTag.RowsAffected() == 0 {
		return repository.ErrNoRows
	}

	return nil
}

func (r *productRepoPostgres) UpdateImage(productID, imageID uuid.UUID) error {
	query := `
		UPDATE products
		SET image_id = $2
		WHERE id = $1
	`

	cmdTag, err := r.pg.Pool.Exec(context.Background(), query, productID, imageID)

	if err != nil {
		return repository.ErrQueryExec
	}

	if cmdTag.RowsAffected() == 0 {
		return repository.ErrNoRows
	}

	return nil
}

func (r *productRepoPostgres) GetByID(id uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	query := `
		SELECT
			id, name, category, price, available_stock,
			last_update_date, supplier_id, image_id
		FROM products
		WHERE id = $1
	`

	err := r.pg.Pool.QueryRow(context.Background(), query, id).Scan(
		&product.ID, &product.Name, &product.Category, &product.Price, &product.AvailableStock,
		&product.LastUpdate, &product.SupplierID, &product.ImageID,
	)

	if err == postgres.ErrNoRows {
		return nil, repository.ErrNoRows
	} else if err != nil {
		return nil, repository.ErrQueryExec
	}

	return &product, nil
}

func (r *productRepoPostgres) GetWithParams(limit, offset int) ([]*entity.Product, error) {
	query := `
		SELECT
			id, name, category, price, available_stock,
			last_update_date, supplier_id, image_id
		FROM products
		LIMIT $1
		OFFSET $2
	`

	rows, err := r.pg.Pool.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, repository.ErrQueryExec
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.ID, &product.Name, &product.Category, &product.Price, &product.AvailableStock,
			&product.LastUpdate, &product.SupplierID, &product.ImageID,
		)
		if err != nil {
			return nil, repository.ErrRowScan
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *productRepoPostgres) RemoveByID(id uuid.UUID) error {
	query := `
		DELETE FROM products
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
