package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type TAGRepository interface {
	Create(auth *model.TAG) error
	GetById(id string) (*model.TAG, error)
	GetByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) bool
	GetDetailByCardNameAndCardNumberAndSetAndDescription(cardName, cardNumber, set, description string) (*model.TAG, error)
	GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) (*model.TAG, error)
	GetByCardNumberAndSetNumber(cardNumber, setNumber string) (*model.TAG, error)
	GetAllCards() ([]*model.TAG, error)
	Update(auth *model.TAG) error
	Delete(id string) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTAGRepository(db *gorm.DB) TAGRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(auth *model.TAG) error {
	result := r.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *tagRepository) GetById(id string) (*model.TAG, error) {
	var auth model.TAG
	result := r.db.First(&auth, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *tagRepository) GetByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) bool {
	var auth model.TAG
	result := r.db.
		Where("cardName = ? AND cardNumber = ? AND setName = ?", cardName, cardNumber, setName).
		First(&auth)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
	}
	return true
}

func (r *tagRepository) GetDetailByCardNameAndCardNumberAndSetAndDescription(cardName, cardNumber, set, description string) (*model.TAG, error) {
	var auth model.TAG
	result := r.db.
		Where("cardName = ? AND cardNumber = ? AND cardSet = ? AND description = ?", cardName, cardNumber, set, description).
		First(&auth)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *tagRepository) GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) (*model.TAG, error) {
	var tag model.TAG
	result := r.db.
		Where("cardName LIKE ? AND cardNumber LIKE ? AND cardSet LIKE ?", "%"+cardName+"%", "%"+cardNumber+"%", "%"+setName+"%").
		First(&tag)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &tag, nil
}

func (r *tagRepository) GetByCardNumberAndSetNumber(cardNumber, setNumber string) (*model.TAG, error) {
	var tag model.TAG
	result := r.db.
		Where("cardNumber = ? AND setNumber = ?", cardNumber, setNumber).
		First(&tag)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &tag, nil
}

func (r *tagRepository) GetAllCards() ([]*model.TAG, error) {
	var tag []*model.TAG
	result := r.db.
		Order("createdDateTime ASC").
		Find(&tag)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return tag, nil
}

func (r *tagRepository) Update(auth *model.TAG) error {
	result := r.db.Save(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *tagRepository) Delete(id string) error {
	result := r.db.Model(&model.TAG{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
