package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type CGCUrlRepository interface {
	Create(auth *model.CGCUrl) error
	GetById(id string) (*model.CGCUrl, error)
	GetByPath(path string) bool
	Update(auth *model.CGCUrl) error
	Delete(id string) error
}

type cgcUrlRepository struct {
	db *gorm.DB
}

func NewCGCUrlRepository(db *gorm.DB) CGCUrlRepository {
	return &cgcUrlRepository{db: db}
}

func (r *cgcUrlRepository) Create(auth *model.CGCUrl) error {
	result := r.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *cgcUrlRepository) GetById(id string) (*model.CGCUrl, error) {
	var auth model.CGCUrl
	result := r.db.First(&auth, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *cgcUrlRepository) GetByPath(path string) bool {
	var url model.CGCUrl
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

func (r *cgcUrlRepository) Update(auth *model.CGCUrl) error {
	result := r.db.Save(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *cgcUrlRepository) Delete(id string) error {
	result := r.db.Model(&model.CGCUrl{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
