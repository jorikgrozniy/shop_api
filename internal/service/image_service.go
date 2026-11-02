package service

import (
	"context"
	"errors"
	"fmt"
	"shop_api/internal/entity"
	"shop_api/internal/repository"
	"shop_api/pkg/postgres"

	"github.com/google/uuid"
)

type ImageService struct {
	imageRepo repository.ImageRepository
}

func NewImageService(imageRepo repository.ImageRepository) *ImageService {
	return &ImageService{
		imageRepo: imageRepo,
	}
}

func (s *ImageService) AddImage(ctx context.Context, q postgres.Querier, image *entity.Image) (uuid.UUID, error) {
	if image == nil {
		return uuid.Nil, fmt.Errorf("add image: %w", ErrNilEntity)
	}

	id, err := s.imageRepo.Save(ctx, q, image)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return id, nil
}

func (s *ImageService) ChangeImage(ctx context.Context, id uuid.UUID, newImage []byte) error {
	err := s.imageRepo.Update(ctx, id, newImage)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return fmt.Errorf("change image %s: %w", id, ErrImageNotFound)
		}
		return fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return nil
}

func (s *ImageService) RemoveImage(ctx context.Context, id uuid.UUID) error {
	err := s.imageRepo.RemoveByID(ctx, id)

	switch {
	case err == nil:
		return nil
	case errors.Is(err, repository.ErrNoRows):
		return fmt.Errorf("remove image %s: %w", id, ErrImageNotFound)
	case errors.Is(err, repository.ErrDependentEntity):
		return fmt.Errorf("remove image %s: %w", id, ErrDependentEntity)
	default:
		return fmt.Errorf("%w: %w", ErrServerInternal, err)
	}
}

func (s *ImageService) GetImage(ctx context.Context, id uuid.UUID) (*entity.Image, error) {
	image, err := s.imageRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, fmt.Errorf("get image %s: %w", id, ErrImageNotFound)
		}
		return nil, fmt.Errorf("%w: %w", ErrServerInternal, err)
	}

	return image, nil
}
