package service

import (
	"pkm/model"
	"pkm/repository"
)

type BGSService struct {
	bgsRepo repository.BGSRepository
}

func NewBGSService(bgsRepo repository.BGSRepository) *BGSService {
	return &BGSService{bgsRepo: bgsRepo}
}

func (s *BGSService) Create(bgs *model.BGS) error {
	return s.bgsRepo.Create(bgs)
}

func (s *BGSService) GetById(id string) (*model.BGS, error) {
	return s.bgsRepo.GetById(id)
}

func (s *BGSService) GetByCardNameAndCardNumberAndSetId(cardName, cardNumber, setId string) bool {
	return s.bgsRepo.GetByCardNameAndCardNumberAndSetId(cardName, cardNumber, setId)
}

func (s *BGSService) GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) (*model.BGS, error) {
	return s.bgsRepo.GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName)
}

func (s *BGSService) GetByCardNumberAndSetNumber(cardNumber, setNumber string) (*model.BGS, error) {
	return s.bgsRepo.GetByCardNumberAndSetNumber(cardNumber, setNumber)
}

func (s *BGSService) GetAllCards() ([]*model.BGS, error) {
	return s.bgsRepo.GetAllCards()
}

func (s *BGSService) Update(bgs *model.BGS) error {
	return s.bgsRepo.Update(bgs)
}

func (s *BGSService) Delete(id string) error {
	return s.bgsRepo.Delete(id)
}
