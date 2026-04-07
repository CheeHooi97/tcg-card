package service

import (
	"pkm/model"
	"pkm/repository"
)

type CGCUrlService struct {
	cgcUrlRepo repository.CGCUrlRepository
}

func NewCGCUrlService(cgcUrlRepo repository.CGCUrlRepository) *CGCUrlService {
	return &CGCUrlService{cgcUrlRepo: cgcUrlRepo}
}

func (s *CGCUrlService) Create(tag *model.CGCUrl) error {
	return s.cgcUrlRepo.Create(tag)
}

func (s *CGCUrlService) GetById(id string) (*model.CGCUrl, error) {
	return s.cgcUrlRepo.GetById(id)
}

func (s *CGCUrlService) GetByPath(path string) bool {
	return s.cgcUrlRepo.GetByPath(path)
}

func (s *CGCUrlService) Update(tag *model.CGCUrl) error {
	return s.cgcUrlRepo.Update(tag)
}

func (s *CGCUrlService) Delete(id string) error {
	return s.cgcUrlRepo.Delete(id)
}
