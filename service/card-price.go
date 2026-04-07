package service

import (
	"pkm/model"
	"pkm/repository"
)

type CardPriceService struct {
	cardPriceRepo repository.CardPriceRepository
}

func NewCardPriceService(cardPriceRepo repository.CardPriceRepository) *CardPriceService {
	return &CardPriceService{cardPriceRepo: cardPriceRepo}
}

func (s *CardPriceService) Create(CardPrice *model.CardPrice) error {
	return s.cardPriceRepo.Create(CardPrice)
}

func (s *CardPriceService) GetById(id string) (*model.CardPrice, error) {
	return s.cardPriceRepo.GetById(id)
}

func (s *CardPriceService) GetByCardPriceNameAndSet(CardPriceName, set string) bool {
	return s.cardPriceRepo.GetByCardPriceNameAndSet(CardPriceName, set)
}

func (s *CardPriceService) Update(CardPrice *model.CardPrice) error {
	return s.cardPriceRepo.Update(CardPrice)
}

func (s *CardPriceService) Delete(id string) error {
	return s.cardPriceRepo.Delete(id)
}
