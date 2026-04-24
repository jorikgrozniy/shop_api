package postgres

import (
	"context"
	"shop_api/internal/dao"
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

func (r *productRepoPostgres) Save(product *dao.Product) error {
	if product == nil {
		return ErrNilEntity
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
		return ErrQueryExec
	}

	return nil
}

func (r *productRepoPostgres) UpdateAvailableStock(id uuid.UUID, value int) error {
	query := `
		UPDATE products
		SET available_stock = available_stock + $2
		WHERE id = $1
	`

	_, err := r.pg.Pool.Exec(context.Background(), query, id, value)

	if err != nil {
		return ErrQueryExec
	}

	return nil
}

func (r *productRepoPostgres) UpdateImage(productID, imageID uuid.UUID) error {
	query := `
		UPDATE products
		SET image_id = $2
		WHERE id = $1
	`

	_, err := r.pg.Pool.Exec(context.Background(), query, productID, imageID)

	if err != nil {
		return ErrQueryExec
	}

	return nil
}

func (r *productRepoPostgres) GetByID(id uuid.UUID) (*dao.Product, error) {
	var product dao.Product
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

	if err != nil {
		return nil, ErrQueryExec
	}

	return &product, nil
}

func (r *productRepoPostgres) GetAll(limit, offset int) ([]*dao.Product, error) {
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
		return nil, ErrQueryExec
	}
	defer rows.Close()

	var products []*dao.Product
	for rows.Next() {
		var product dao.Product
		err := rows.Scan(
			&product.ID, &product.Name, &product.Category, &product.Price, &product.AvailableStock,
			&product.LastUpdate, &product.SupplierID, &product.ImageID,
		)
		if err != nil {
			return nil, ErrRowScan
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

	_, err := r.pg.Pool.Exec(context.Background(), query, id)

	if err != nil {
		return ErrQueryExec
	}

	return nil
}
