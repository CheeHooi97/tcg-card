package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type BGSLoggingRepository interface {
	Create(auth *model.BGSLogging) error
	GetById(id string) (*model.BGSLogging, error)
	Update(auth *model.BGSLogging) error
	Delete(id string) error
}

type bgsLoggingRepository struct {
	db *gorm.DB
}

func NewBGSLoggingRepository(db *gorm.DB) BGSLoggingRepository {
	return &bgsLoggingRepository{db: db}
}

func (r *bgsLoggingRepository) Create(auth *model.BGSLogging) error {
	result := r.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *bgsLoggingRepository) GetById(id string) (*model.BGSLogging, error) {
	var auth model.BGSLogging
	result := r.db.First(&auth, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *bgsLoggingRepository) Update(auth *model.BGSLogging) error {
	result := r.db.Save(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *bgsLoggingRepository) Delete(id string) error {
	result := r.db.Model(&model.BGSLogging{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
