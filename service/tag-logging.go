package service

import (
	"pkm/model"
	"pkm/repository"
)

type TAGLoggingService struct {
	tagLoggingRepo repository.TAGLoggingRepository
}

func NewTAGLoggingService(tagLoggingRepo repository.TAGLoggingRepository) *TAGLoggingService {
	return &TAGLoggingService{tagLoggingRepo: tagLoggingRepo}
}

func (s *TAGLoggingService) Create(tag *model.TAGLogging) error {
	return s.tagLoggingRepo.Create(tag)
}

func (s *TAGLoggingService) GetById(id string) (*model.TAGLogging, error) {
	return s.tagLoggingRepo.GetById(id)
}

func (s *TAGLoggingService) Update(tag *model.TAGLogging) error {
	return s.tagLoggingRepo.Update(tag)
}

func (s *TAGLoggingService) Delete(id string) error {
	return s.tagLoggingRepo.Delete(id)
}
