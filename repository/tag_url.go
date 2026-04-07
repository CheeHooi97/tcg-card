package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type TAGUrlRepository interface {
	Create(url *model.TAGUrl) error
	GetById(id string) (*model.TAGUrl, error)
	GetByPath(path string) bool
	Update(url *model.TAGUrl) error
	Delete(id string) error
}

type tagUrlRepository struct {
	db *gorm.DB
}

func NewTAGUrlRepository(db *gorm.DB) TAGUrlRepository {
	return &tagUrlRepository{db: db}
}

func (r *tagUrlRepository) Create(url *model.TAGUrl) error {
	result := r.db.Create(url)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *tagUrlRepository) GetById(id string) (*model.TAGUrl, error) {
	var url model.TAGUrl
	result := r.db.First(&url, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &url, nil
}

func (r *tagUrlRepository) GetByPath(path string) bool {
	var url model.TAGUrl
	result := r.db.
		Where("url = ?", path).
		First(&url)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
	}
	return true
}

func (r *tagUrlRepository) Update(url *model.TAGUrl) error {
	result := r.db.Save(url)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *tagUrlRepository) Delete(id string) error {
	result := r.db.Model(&model.TAGUrl{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
