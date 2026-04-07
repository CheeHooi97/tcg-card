package service

import (
	"pkm/model"
	"pkm/repository"
)

type BGSLoggingService struct {
	bgsLoggingRepo repository.BGSLoggingRepository
}

func NewBGSLoggingService(bgsLoggingRepo repository.BGSLoggingRepository) *BGSLoggingService {
	return &BGSLoggingService{bgsLoggingRepo: bgsLoggingRepo}
}

func (s *BGSLoggingService) Create(bgs *model.BGSLogging) error {
	return s.bgsLoggingRepo.Create(bgs)
}

func (s *BGSLoggingService) GetById(id string) (*model.BGSLogging, error) {
	return s.bgsLoggingRepo.GetById(id)
}

func (s *BGSLoggingService) Update(bgs *model.BGSLogging) error {
	return s.bgsLoggingRepo.Update(bgs)
}

func (s *BGSLoggingService) Delete(id string) error {
	return s.bgsLoggingRepo.Delete(id)
}
