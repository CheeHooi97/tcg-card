package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type BGSUrlRepository interface {
	Create(auth *model.BGSUrl) error
	GetById(id string) (*model.BGSUrl, error)
	GetByPath(path string) bool
	Update(auth *model.BGSUrl) error
	Delete(id string) error
}

type bgsUrlRepository struct {
	db *gorm.DB
}

func NewBGSUrlRepository(db *gorm.DB) BGSUrlRepository {
	return &bgsUrlRepository{db: db}
}

func (r *bgsUrlRepository) Create(auth *model.BGSUrl) error {
	result := r.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *bgsUrlRepository) GetById(id string) (*model.BGSUrl, error) {
	var auth model.BGSUrl
	result := r.db.First(&auth, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *bgsUrlRepository) GetByPath(path string) bool {
	var url model.BGSUrl
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

func (r *bgsUrlRepository) Update(auth *model.BGSUrl) error {
	result := r.db.Save(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *bgsUrlRepository) Delete(id string) error {
	result := r.db.Model(&model.BGSUrl{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
