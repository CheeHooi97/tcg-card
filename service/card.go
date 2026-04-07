package service

import (
	"pkm/model"
	"pkm/repository"
)

type CardService struct {
	cardRepo repository.CardRepository
}

func NewCardService(cardRepo repository.CardRepository) *CardService {
	return &CardService{cardRepo: cardRepo}
}

func (s *CardService) Create(card *model.Card) error {
	return s.cardRepo.Create(card)
}

func (s *CardService) GetById(id string) (*model.Card, error) {
	return s.cardRepo.GetById(id)
}

func (s *CardService) GetByCardNameAndSet(cardName, set string) bool {
	return s.cardRepo.GetByCardNameAndSet(cardName, set)
}

func (s *CardService) SearchCardKeywords(keyword string) ([]*model.Card, error) {
	return s.cardRepo.SearchCardKeywords(keyword)
}

func (s *CardService) SearchCardBySort(sortType, userId string) ([]*model.Card, error) {
	return s.cardRepo.SearchCardBySort(sortType, userId)
}

func (s *CardService) GetAllCards() ([]*model.Card, error) {
	return s.cardRepo.GetAllCards()
}

func (s *CardService) Update(card *model.Card) error {
	return s.cardRepo.Update(card)
}

func (s *CardService) Delete(id string) error {
	return s.cardRepo.Delete(id)
}
