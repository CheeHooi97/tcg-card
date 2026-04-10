package repository

import (
	"errors"
	"fmt"
	"pkm/model"
	"time"

	"gorm.io/gorm"
)

type PSARepository interface {
	Create(auth *model.PSA) error
	GetById(id string) (*model.PSA, error)
	CheckBySpecId(specId string) bool
	CheckPopulation(pop, description string) ([]*model.PSA, error)
	GetSpecIDs() ([]*model.PSA, error)
	GetByCardNameAndCardNumberAndDescriptionAndSet(cardName, cardNumber, description, set string) bool
	GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) (*model.PSA, error)
	GetDetailByCardNameAndCardNumberAndDescription(cardName, cardNumber, description string) (*model.PSA, error)
	GetByCardNumberAndSetNumber(cardNumber, setNumber string) (*model.PSA, error)
	Update(auth *model.PSA) error
	Delete(id string) error
}

type psaRepository struct {
	db *gorm.DB
}

func NewPSARepository(db *gorm.DB) PSARepository {
	return &psaRepository{db: db}
}

func (r *psaRepository) Create(auth *model.PSA) error {
	result := r.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *psaRepository) GetById(id string) (*model.PSA, error) {
	var auth model.PSA
	result := r.db.First(&auth, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *psaRepository) CheckBySpecId(specId string) bool {
	var auth model.PSA
	result := r.db.
		Where("specId = ?", specId).
		First(&auth)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
	}

	return true
}

func (r *psaRepository) CheckPopulation(pop, description string) ([]*model.PSA, error) {
	var auth []*model.PSA
	fmt.Println("a:", "%"+description+"%")
	result := r.db.
		Where("total <= ? AND description LIKE ?", pop, "%"+description+"%").
		Find(&auth)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}

	return auth, nil
}

func (r *psaRepository) GetSpecIDs() ([]*model.PSA, error) {
	var auth []*model.PSA

	startOfToday := time.Now().UTC().Truncate(24 * time.Hour)

	result := r.db.
		Where("updated_date_time < ?", startOfToday).
		Find(&auth)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}

	return auth, nil
}

func (r *psaRepository) GetByCardNameAndCardNumberAndDescriptionAndSet(cardName, cardNumber, description, setName string) bool {
	var auth model.PSA
	result := r.db.
		Where("card_name = ? AND card_number = ? AND description = ? AND set_name = ?", cardName, cardNumber, description, setName).
		First(&auth)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
	}
	return true
}

func (r *psaRepository) GetDetailByCardNameAndCardNumberAndSetName(cardName, cardNumber, setName string) (*model.PSA, error) {
	var psa model.PSA
	result := r.db.
		Where("card_name LIKE ? AND card_number LIKE ? AND set_name LIKE ?", "%"+cardName+"%", "%"+cardNumber+"%", "%"+setName+"%").
		First(&psa)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &psa, nil
}

func (r *psaRepository) GetDetailByCardNameAndCardNumberAndDescription(cardName, cardNumber, description string) (*model.PSA, error) {
	var auth model.PSA
	result := r.db.
		Where("card_name = ? AND card_number = ? AND description = ?", cardName, cardNumber, description).
		First(&auth)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *psaRepository) GetByCardNumberAndSetNumber(cardNumber, setNumber string) (*model.PSA, error) {
	var auth model.PSA
	result := r.db.
		Where("card_number = ? AND set_number = ?", cardNumber, setNumber).
		First(&auth)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *psaRepository) Update(auth *model.PSA) error {
	result := r.db.Save(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *psaRepository) Delete(id string) error {
	result := r.db.Model(&model.PSA{}).Where("id = ?", id).Update("status", false)
	return result.Error
}

