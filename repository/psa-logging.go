package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type PSALoggingRepository interface {
	Create(auth *model.PSALogging) error
	GetById(id string) (*model.PSALogging, error)
	Update(auth *model.PSALogging) error
	Delete(id string) error
}

type psaLoggingRepository struct {
	db *gorm.DB
}

func NewPSALoggingRepository(db *gorm.DB) PSALoggingRepository {
	return &psaLoggingRepository{db: db}
}

func (r *psaLoggingRepository) Create(auth *model.PSALogging) error {
	result := r.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *psaLoggingRepository) GetById(id string) (*model.PSALogging, error) {
	var auth model.PSALogging
	result := r.db.First(&auth, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *psaLoggingRepository) Update(auth *model.PSALogging) error {
	result := r.db.Save(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *psaLoggingRepository) Delete(id string) error {
	result := r.db.Model(&model.PSALogging{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
