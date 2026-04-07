package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type PSAUrlRepository interface {
	Create(url *model.PSAUrl) error
	GetById(id string) (*model.PSAUrl, error)
	GetByPath(path string) bool
	Update(url *model.PSAUrl) error
	Delete(id string) error
}

type psaUrlRepository struct {
	db *gorm.DB
}

func NewPSAUrlRepository(db *gorm.DB) PSAUrlRepository {
	return &psaUrlRepository{db: db}
}

func (r *psaUrlRepository) Create(url *model.PSAUrl) error {
	result := r.db.Create(url)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *psaUrlRepository) GetById(id string) (*model.PSAUrl, error) {
	var url model.PSAUrl
	result := r.db.First(&url, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &url, nil
}

func (r *psaUrlRepository) GetByPath(path string) bool {
	var url model.PSAUrl
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

func (r *psaUrlRepository) Update(url *model.PSAUrl) error {
	result := r.db.Save(url)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *psaUrlRepository) Delete(id string) error {
	result := r.db.Model(&model.PSAUrl{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
