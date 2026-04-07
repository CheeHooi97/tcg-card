package service

import (
	"pkm/model"
	"pkm/repository"
)

type TAGService struct {
	tagRepo repository.TAGRepository
}

func NewTAGService(tagRepo repository.TAGRepository) *TAGService {
	return &TAGService{tagRepo: tagRepo}
}

func (s *TAGService) Create(tag *model.TAG) error {
	return s.tagRepo.Create(tag)
}

func (s *TAGService) GetById(id string) (*model.TAG, error) {
	return s.tagRepo.GetById(id)
}

func (s *TAGService) GetByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) bool {
	return s.tagRepo.GetByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName)
}

func (s *TAGService) GetDetailByCardNameAndCardNumberAndSetAndDescription(cardName, cardNumber, set, description string) (*model.TAG, error) {
	return s.tagRepo.GetDetailByCardNameAndCardNumberAndSetAndDescription(cardName, cardNumber, set, description)
}

func (s *TAGService) GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) (*model.TAG, error) {
	return s.tagRepo.GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName)
}

func (s *TAGService) GetByCardNumberAndSetNumber(cardNumber, setNumber string) (*model.TAG, error) {
	return s.tagRepo.GetByCardNumberAndSetNumber(cardNumber, setNumber)
}

func (s *TAGService) GetAllCards() ([]*model.TAG, error) {
	return s.tagRepo.GetAllCards()
}

func (s *TAGService) Update(tag *model.TAG) error {
	return s.tagRepo.Update(tag)
}

func (s *TAGService) Delete(id string) error {
	return s.tagRepo.Delete(id)
}
