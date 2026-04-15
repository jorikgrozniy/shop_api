package service

import (
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

func (s *ImageService) AddImage(image []byte) (uuid.UUID, error) {
	id, err := s.imageRepo.Save(image)
	if err != nil {
		return uuid.Nil, ErrServerInternal
	}

	return id, nil
}

func (s *ImageService) ChangeImage(id uuid.UUID, newImage []byte) error {
	return s.imageRepo.Update(id, newImage)
}

func (s *ImageService) RemoveImage(id uuid.UUID) error {
	return s.imageRepo.RemoveByID(id)
}

func (s *ImageService) GetImage(id uuid.UUID) ([]byte, error) {
	return s.imageRepo.GetByID(id)
}
