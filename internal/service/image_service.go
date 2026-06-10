package service

import (
	"shop_api/internal/entity"
	"shop_api/internal/repository"

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

func (s *ImageService) AddImage(image *entity.Image) (uuid.UUID, error) {
	id, err := s.imageRepo.Save(image)
	if err != nil {
		return uuid.Nil, ErrServerInternal
	}

	return id, nil
}

func (s *ImageService) ChangeImage(id uuid.UUID, newImage []byte) error {
	if err := s.imageRepo.Update(id, newImage); err == repository.ErrNoRows {
		return ErrImageNotFound
	} else if err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ImageService) RemoveImage(id uuid.UUID) error {
	err := s.imageRepo.RemoveByID(id)

	switch err {
	case repository.ErrNoRows:
		return ErrImageNotFound
	case repository.ErrDependentEntity:
		return ErrDependentEntity
	}

	if err != nil {
		return ErrServerInternal
	}

	return nil
}

func (s *ImageService) GetImage(id uuid.UUID) (*entity.Image, error) {
	image, err := s.imageRepo.GetByID(id)

	if err == repository.ErrNoRows {
		return nil, ErrImageNotFound
	} else if err != nil {
		return nil, ErrServerInternal
	}

	return image, nil
}
