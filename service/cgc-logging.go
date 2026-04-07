package service

import (
	"pkm/model"
	"pkm/repository"
)

type CGCLoggingService struct {
	cgcLoggingRepo repository.CGCLoggingRepository
}

func NewCGCLoggingService(cgcLoggingRepo repository.CGCLoggingRepository) *CGCLoggingService {
	return &CGCLoggingService{cgcLoggingRepo: cgcLoggingRepo}
}

func (s *CGCLoggingService) Create(bgs *model.CGCLogging) error {
	return s.cgcLoggingRepo.Create(bgs)
}

func (s *CGCLoggingService) GetById(id string) (*model.CGCLogging, error) {
	return s.cgcLoggingRepo.GetById(id)
}

func (s *CGCLoggingService) Update(bgs *model.CGCLogging) error {
	return s.cgcLoggingRepo.Update(bgs)
}

func (s *CGCLoggingService) Delete(id string) error {
	return s.cgcLoggingRepo.Delete(id)
}
