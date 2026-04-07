package service

import (
	"pkm/model"
	"pkm/repository"
)

type TAGUrlService struct {
	tagUrlRepo repository.TAGUrlRepository
}

func NewTAGUrlService(tagUrlRepo repository.TAGUrlRepository) *TAGUrlService {
	return &TAGUrlService{tagUrlRepo: tagUrlRepo}
}

func (s *TAGUrlService) Create(tag *model.TAGUrl) error {
	return s.tagUrlRepo.Create(tag)
}

func (s *TAGUrlService) GetById(id string) (*model.TAGUrl, error) {
	return s.tagUrlRepo.GetById(id)
}

func (s *TAGUrlService) GetByPath(path string) bool {
	return s.tagUrlRepo.GetByPath(path)
}

func (s *TAGUrlService) Update(tag *model.TAGUrl) error {
	return s.tagUrlRepo.Update(tag)
}

func (s *TAGUrlService) Delete(id string) error {
	return s.tagUrlRepo.Delete(id)
}
