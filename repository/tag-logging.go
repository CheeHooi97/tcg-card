package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type TAGLoggingRepository interface {
	Create(auth *model.TAGLogging) error
	GetById(id string) (*model.TAGLogging, error)
	Update(auth *model.TAGLogging) error
	Delete(id string) error
}

type tagLoggingRepository struct {
	db *gorm.DB
}

func NewTAGLoggingRepository(db *gorm.DB) TAGLoggingRepository {
	return &tagLoggingRepository{db: db}
}

func (r *tagLoggingRepository) Create(auth *model.TAGLogging) error {
	result := r.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *tagLoggingRepository) GetById(id string) (*model.TAGLogging, error) {
	var auth model.TAGLogging
	result := r.db.First(&auth, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *tagLoggingRepository) Update(auth *model.TAGLogging) error {
	result := r.db.Save(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *tagLoggingRepository) Delete(id string) error {
	result := r.db.Model(&model.TAGLogging{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
