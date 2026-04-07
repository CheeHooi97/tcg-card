package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type CardPriceRepository interface {
	Create(CardPrice *model.CardPrice) error
	GetById(id string) (*model.CardPrice, error)
	GetByCardPriceNameAndSet(CardPriceName, set string) bool
	Update(CardPrice *model.CardPrice) error
	Delete(id string) error
}

type cardPriceRepository struct {
	db *gorm.DB
}

func NewCardPriceRepository(db *gorm.DB) CardPriceRepository {
	return &cardPriceRepository{db: db}
}

func (r *cardPriceRepository) Create(CardPrice *model.CardPrice) error {
	result := r.db.Create(CardPrice)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *cardPriceRepository) GetById(id string) (*model.CardPrice, error) {
	var CardPrice model.CardPrice
	result := r.db.First(&CardPrice, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &CardPrice, nil
}

func (r *cardPriceRepository) GetByCardPriceNameAndSet(CardPriceName, set string) bool {
	var auth model.CardPrice
	result := r.db.
		Where("name = ? AND set = ?", CardPriceName, set).
		First(&auth)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
	}
	return true
}

func (r *cardPriceRepository) Update(CardPrice *model.CardPrice) error {
	result := r.db.Save(CardPrice)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *cardPriceRepository) Delete(id string) error {
	result := r.db.Model(&model.CardPrice{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
