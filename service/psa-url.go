package service

import (
	"pkm/model"
	"pkm/repository"
)

type PSAUrlService struct {
	psaUrlRepo repository.PSAUrlRepository
}

func NewPSAUrlService(psaUrlRepo repository.PSAUrlRepository) *PSAUrlService {
	return &PSAUrlService{psaUrlRepo: psaUrlRepo}
}

func (s *PSAUrlService) Create(tag *model.PSAUrl) error {
	return s.psaUrlRepo.Create(tag)
}

func (s *PSAUrlService) GetById(id string) (*model.PSAUrl, error) {
	return s.psaUrlRepo.GetById(id)
}

func (s *PSAUrlService) GetByPath(path string) bool {
	return s.psaUrlRepo.GetByPath(path)
}

func (s *PSAUrlService) Update(tag *model.PSAUrl) error {
	return s.psaUrlRepo.Update(tag)
}

func (s *PSAUrlService) Delete(id string) error {
	return s.psaUrlRepo.Delete(id)
}
