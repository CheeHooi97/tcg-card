package service

import (
	"pkm/model"
	"pkm/repository"
)

type PSAService struct {
	psaRepo repository.PSARepository
}

func NewPSAService(psaRepo repository.PSARepository) *PSAService {
	return &PSAService{psaRepo: psaRepo}
}

func (s *PSAService) Create(authenticator *model.PSA) error {
	return s.psaRepo.Create(authenticator)
}

func (s *PSAService) GetById(id string) (*model.PSA, error) {
	return s.psaRepo.GetById(id)
}

func (s *PSAService) CheckBySpecId(specId string) bool {
	return s.psaRepo.CheckBySpecId(specId)
}

func (s *PSAService) CheckPopulation(pop, description string) ([]*model.PSA, error) {
	return s.psaRepo.CheckPopulation(pop, description)
}

func (s *PSAService) GetSpecIDs() ([]*model.PSA, error) {
	return s.psaRepo.GetSpecIDs()
}

func (s *PSAService) GetByCardNameAndCardNumberAndDescriptionAndSet(cardName, cardNumber, description, setName string) bool {
	return s.psaRepo.GetByCardNameAndCardNumberAndDescriptionAndSet(cardName, cardNumber, description, setName)
}

func (s *PSAService) GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) (*model.PSA, error) {
	return s.psaRepo.GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName)
}

func (s *PSAService) GetDetailByCardNameAndCardNumberAndDescription(cardName, cardNumber, description string) (*model.PSA, error) {
	return s.psaRepo.GetDetailByCardNameAndCardNumberAndDescription(cardName, cardNumber, description)
}

func (s *PSAService) GetByCardNumberAndSetNumber(cardNumber, setNumber string) (*model.PSA, error) {
	return s.psaRepo.GetByCardNumberAndSetNumber(cardNumber, setNumber)
}

func (s *PSAService) Update(authenticator *model.PSA) error {
	return s.psaRepo.Update(authenticator)
}

func (s *PSAService) Delete(id string) error {
	return s.psaRepo.Delete(id)
}
