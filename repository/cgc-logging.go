package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type CGCLoggingRepository interface {
	Create(auth *model.CGCLogging) error
	GetById(id string) (*model.CGCLogging, error)
	Update(auth *model.CGCLogging) error
	Delete(id string) error
}

type cgcLoggingRepository struct {
	db *gorm.DB
}

func NewCGCLoggingRepository(db *gorm.DB) CGCLoggingRepository {
	return &cgcLoggingRepository{db: db}
}

func (r *cgcLoggingRepository) Create(auth *model.CGCLogging) error {
	result := r.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *cgcLoggingRepository) GetById(id string) (*model.CGCLogging, error) {
	var auth model.CGCLogging
	result := r.db.First(&auth, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *cgcLoggingRepository) Update(auth *model.CGCLogging) error {
	result := r.db.Save(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *cgcLoggingRepository) Delete(id string) error {
	result := r.db.Model(&model.CGCLogging{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
