package service

import (
	"pkm/model"
	"pkm/repository"
)

type CGCService struct {
	cgcRepo repository.CGCRepository
}

func NewCGCService(cgcRepo repository.CGCRepository) *CGCService {
	return &CGCService{cgcRepo: cgcRepo}
}

func (s *CGCService) Create(bgs *model.CGC) error {
	return s.cgcRepo.Create(bgs)
}

func (s *CGCService) GetById(id string) (*model.CGC, error) {
	return s.cgcRepo.GetById(id)
}

func (s *CGCService) CheckCardNameAndCardNumberAndSetNameAndRarity(cardName, cardNumber, setName, rarity string) bool {
	return s.cgcRepo.CheckCardNameAndCardNumberAndSetNameAndRarity(cardName, cardNumber, setName, rarity)
}

func (s *CGCService) GetByCardNameAndCardNumberAndSetNameAndRarity(cardName, cardNumber, setName, rarity string) (*model.CGC, error) {
	return s.cgcRepo.GetByCardNameAndCardNumberAndSetNameAndRarity(cardName, cardNumber, setName, rarity)
}

func (s *CGCService) GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) (*model.CGC, error) {
	return s.cgcRepo.GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName)
}

func (s *CGCService) GetByCardNumberAndSetNumber(cardNumber, setNumber string) (*model.CGC, error) {
	return s.cgcRepo.GetByCardNumberAndSetNumber(cardNumber, setNumber)
}

func (s *CGCService) GetAllCards() ([]*model.CGC, error) {
	return s.cgcRepo.GetAllCards()
}

func (s *CGCService) Update(bgs *model.CGC) error {
	return s.cgcRepo.Update(bgs)
}

func (s *CGCService) Delete(id string) error {
	return s.cgcRepo.Delete(id)
}
