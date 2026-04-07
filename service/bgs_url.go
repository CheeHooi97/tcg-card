package service

import (
	"pkm/model"
	"pkm/repository"
)

type BGSUrlService struct {
	bgsUrlRepo repository.BGSUrlRepository
}

func NewBGSUrlService(bgsUrlRepo repository.BGSUrlRepository) *BGSUrlService {
	return &BGSUrlService{bgsUrlRepo: bgsUrlRepo}
}

func (s *BGSUrlService) Create(tag *model.BGSUrl) error {
	return s.bgsUrlRepo.Create(tag)
}

func (s *BGSUrlService) GetById(id string) (*model.BGSUrl, error) {
	return s.bgsUrlRepo.GetById(id)
}

func (s *BGSUrlService) GetByPath(path string) bool {
	return s.bgsUrlRepo.GetByPath(path)
}

func (s *BGSUrlService) Update(tag *model.BGSUrl) error {
	return s.bgsUrlRepo.Update(tag)
}

func (s *BGSUrlService) Delete(id string) error {
	return s.bgsUrlRepo.Delete(id)
}
