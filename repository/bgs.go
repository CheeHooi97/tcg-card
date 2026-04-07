package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type BGSRepository interface {
	Create(auth *model.BGS) error
	GetById(id string) (*model.BGS, error)
	GetByCardNameAndCardNumberAndSetId(cardName, cardNumber, setId string) bool
	GetDetailByCardNameAndCardNumberAndSetId(cardName, cardNumber, setId string) (*model.BGS, error)
	GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) (*model.BGS, error)
	GetByCardNumberAndSetNumber(cardNumber, setNumber string) (*model.BGS, error)
	GetAllCards() ([]*model.BGS, error)
	Update(auth *model.BGS) error
	Delete(id string) error
}

type bgsRepository struct {
	db *gorm.DB
}

func NewBGSRepository(db *gorm.DB) BGSRepository {
	return &bgsRepository{db: db}
}

func (r *bgsRepository) Create(auth *model.BGS) error {
	result := r.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *bgsRepository) GetById(id string) (*model.BGS, error) {
	var auth model.BGS
	result := r.db.First(&auth, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *bgsRepository) GetByCardNameAndCardNumberAndSetId(cardName, cardNumber, setId string) bool {
	var auth model.BGS
	result := r.db.
		Where("cardName = ? AND cardNumber = ? AND setId = ?", cardName, cardNumber, setId).
		First(&auth)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
	}
	return true
}

func (r *bgsRepository) GetDetailByCardNameAndCardNumberAndSetId(cardName, cardNumber, setId string) (*model.BGS, error) {
	var bgs model.BGS
	result := r.db.
		Where("cardName = ? AND cardNumber = ? AND setId = ?", cardName, cardNumber, setId).
		First(&bgs)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &bgs, nil
}

func (r *bgsRepository) GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) (*model.BGS, error) {
	var bgs model.BGS
	result := r.db.
		Where("cardName LIKE ? AND cardNumber LIKE ? AND description LIKE ?", "%"+cardName+"%", "%"+cardNumber+"%", "%"+setName+"%").
		First(&bgs)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &bgs, nil
}

func (r *bgsRepository) GetByCardNumberAndSetNumber(cardNumber, setNumber string) (*model.BGS, error) {
	var bgs model.BGS
	result := r.db.
		Where("cardNumber = ? AND setNumber = ?", cardNumber, setNumber).
		First(&bgs)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &bgs, nil
}

func (r *bgsRepository) GetAllCards() ([]*model.BGS, error) {
	var bgs []*model.BGS
	result := r.db.
		Order("createdDateTime ASC").
		Find(&bgs)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return bgs, nil
}

func (r *bgsRepository) Update(auth *model.BGS) error {
	result := r.db.Save(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *bgsRepository) Delete(id string) error {
	result := r.db.Model(&model.BGS{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
