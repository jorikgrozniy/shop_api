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

type imageRepoPostgres struct {
	pg *postgres.Postgres
}

func NewImageRepoPostgres(pg *postgres.Postgres) repository.ImageRepository {
	return &imageRepoPostgres{
		pg: pg,
	}
}

func (r *imageRepoPostgres) Save(ctx context.Context, q postgres.Querier, image *entity.Image) (uuid.UUID, error) {
	if image == nil {
		return uuid.Nil, fmt.Errorf("save image: %w", repository.ErrNilEntity)
	}

	query := `
		INSERT INTO images (image)
		VALUES ($1)
		RETURNING id
	`

	var id uuid.UUID
	err := q.QueryRow(ctx, query, image.Image).Scan(&id)

	if errors.Is(err, postgres.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("save image: %w", repository.ErrNoRows)
	} else if err != nil {
		return uuid.Nil, fmt.Errorf("save image: %w: %w", repository.ErrQueryExec, err)
	}

	return id, nil
}

func (r *imageRepoPostgres) RemoveByID(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM images
		WHERE id = $1
	`

	cmdTag, err := r.pg.Pool.Exec(ctx, query, id)

	if isFKViolation(err) {
		return fmt.Errorf("remove image %s: %w: %w", id, repository.ErrDependentEntity, err)
	} else if err != nil {
		return fmt.Errorf("remove image %s: %w: %w", id, repository.ErrQueryExec, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("remove image %s: %w", id, repository.ErrNoRows)
	}

	return nil
}

func (r *imageRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.Image, error) {
	var image entity.Image
	query := `
		SELECT image
		FROM images
		WHERE id = $1
	`

	err := r.pg.Pool.QueryRow(ctx, query, id).Scan(&image.Image)

	if errors.Is(err, postgres.ErrNoRows) {
		return nil, fmt.Errorf("get image %s: %w", id, repository.ErrNoRows)
	} else if err != nil {
		return nil, fmt.Errorf("get image %s: %w: %w", id, repository.ErrQueryExec, err)
	}

	return &image, nil
}

func (r *imageRepoPostgres) Update(ctx context.Context, id uuid.UUID, newImage []byte) error {
	query := `
		UPDATE images
		SET image = $2
		WHERE id = $1
	`

	cmdTag, err := r.pg.Pool.Exec(ctx, query, id, newImage)

	if err != nil {
		return fmt.Errorf("update image %s: %w: %w", id, repository.ErrQueryExec, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("update image %s: %w", id, repository.ErrNoRows)
	}

	return nil
}
