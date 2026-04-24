package postgres

import (
	"context"
	"shop_api/internal/repository"
	"shop_api/pkg/postgres"

	"github.com/google/uuid"
)

type imageRepoPostgres struct {
	pg *postgres.Postgres
}

func NewImageRepoPostgres(pg *postgres.Postgres) repository.ImageRepository {
	return &imageRepoPostgres{
		pg: pg,
	}
}

func (r *imageRepoPostgres) Save(image []byte) (uuid.UUID, error) {
	if image == nil {
		return uuid.Nil, ErrNilEntity
	}

	query := `
		INSERT INTO images (image)
		VALUES ($1)
		RETURNING id
	`

	var id uuid.UUID
	err := r.pg.Pool.QueryRow(context.Background(), query, image).Scan(&id)

	if err != nil {
		return uuid.Nil, ErrQueryExec
	}

	return id, nil
}

func (r *imageRepoPostgres) RemoveByID(id uuid.UUID) error {
	query := `
		DELETE FROM images
		WHERE id = $1
	`

	_, err := r.pg.Pool.Exec(context.Background(), query, id)

	if err != nil {
		return ErrQueryExec
	}

	return nil
}

func (r *imageRepoPostgres) GetByID(id uuid.UUID) ([]byte, error) {
	var image []byte
	query := `
		SELECT image
		FROM images
		WHERE id = $1
	`

	err := r.pg.Pool.QueryRow(context.Background(), query, id).Scan(&image)

	if err != nil {
		return nil, ErrQueryExec
	}

	return image, nil
}

func (r *imageRepoPostgres) Update(id uuid.UUID, newImage []byte) error {
	query := `
		UPDATE images
		SET image = $2
		WHERE id = $1
	`

	_, err := r.pg.Pool.Exec(context.Background(), query, id, newImage)

	if err != nil {
		return ErrQueryExec
	}

	return nil
}
