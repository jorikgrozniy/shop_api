package postgres

import (
	"context"
	"shop_api/internal/entity"
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

func (r *imageRepoPostgres) Save(image *entity.Image) (uuid.UUID, error) {
	if image == nil {
		return uuid.Nil, repository.ErrNilEntity
	}

	query := `
		INSERT INTO images (image)
		VALUES ($1)
		RETURNING id
	`

	var id uuid.UUID
	err := r.pg.Pool.QueryRow(context.Background(), query, image.Image).Scan(&id)

	if err == postgres.ErrNoRows {
		return uuid.Nil, repository.ErrNoRows
	} else if err != nil {
		return uuid.Nil, repository.ErrQueryExec
	}

	return id, nil
}

func (r *imageRepoPostgres) RemoveByID(id uuid.UUID) error {
	query := `
		DELETE FROM images
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

func (r *imageRepoPostgres) GetByID(id uuid.UUID) (*entity.Image, error) {
	var image entity.Image
	query := `
		SELECT image
		FROM images
		WHERE id = $1
	`

	err := r.pg.Pool.QueryRow(context.Background(), query, id).Scan(&image.Image)

	if err == postgres.ErrNoRows {
		return nil, repository.ErrNoRows
	} else if err != nil {
		return nil, repository.ErrQueryExec
	}

	return &image, nil
}

func (r *imageRepoPostgres) Update(id uuid.UUID, newImage []byte) error {
	query := `
		UPDATE images
		SET image = $2
		WHERE id = $1
	`

	cmdTag, err := r.pg.Pool.Exec(context.Background(), query, id, newImage)

	if err != nil {
		return repository.ErrQueryExec
	}

	if cmdTag.RowsAffected() == 0 {
		return repository.ErrNoRows
	}

	return nil
}
