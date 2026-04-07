package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type SetRepository interface {
	Create(set *model.Set) error
	GetById(id string) (*model.Set, error)
	Update(set *model.Set) error
	Delete(id string) error
}

type setRepository struct {
	db *gorm.DB
}

func NewSetRepository(db *gorm.DB) SetRepository {
	return &setRepository{db: db}
}

func (r *setRepository) Create(set *model.Set) error {
	result := r.db.Create(set)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *setRepository) GetById(id string) (*model.Set, error) {
	var set model.Set
	result := r.db.First(&set, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &set, nil
}

func (r *setRepository) Update(set *model.Set) error {
	result := r.db.Save(set)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *setRepository) Delete(id string) error {
	result := r.db.Model(&model.Set{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
