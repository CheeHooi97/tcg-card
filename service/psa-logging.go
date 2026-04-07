package service

import (
	"pkm/model"
	"pkm/repository"
)

type PSALoggingService struct {
	psaLoggingRepo repository.PSALoggingRepository
}

func NewPSALoggingService(psaLoggingRepo repository.PSALoggingRepository) *PSALoggingService {
	return &PSALoggingService{psaLoggingRepo: psaLoggingRepo}
}

func (s *PSALoggingService) Create(authenticator *model.PSALogging) error {
	return s.psaLoggingRepo.Create(authenticator)
}

func (s *PSALoggingService) GetById(id string) (*model.PSALogging, error) {
	return s.psaLoggingRepo.GetById(id)
}

func (s *PSALoggingService) Update(authenticator *model.PSALogging) error {
	return s.psaLoggingRepo.Update(authenticator)
}

func (s *PSALoggingService) Delete(id string) error {
	return s.psaLoggingRepo.Delete(id)
}
