package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type UserCardSearchLogRepository interface {
	Create(auth *model.UserCardSearchLog) error
	GetById(id string) (*model.UserCardSearchLog, error)
	GetAllCards() ([]*model.UserCardSearchLog, error)
	Update(auth *model.UserCardSearchLog) error
	Delete(id string) error
}

type userCardSearchLogRepository struct {
	db *gorm.DB
}

func NewUserCardSearchLogRepository(db *gorm.DB) UserCardSearchLogRepository {
	return &userCardSearchLogRepository{db: db}
}

func (r *userCardSearchLogRepository) Create(auth *model.UserCardSearchLog) error {
	result := r.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userCardSearchLogRepository) GetById(id string) (*model.UserCardSearchLog, error) {
	var auth model.UserCardSearchLog
	result := r.db.First(&auth, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *userCardSearchLogRepository) GetAllCards() ([]*model.UserCardSearchLog, error) {
	var UserCardSearchLog []*model.UserCardSearchLog
	result := r.db.
		Order("created_date_time ASC").
		Find(&UserCardSearchLog)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return UserCardSearchLog, nil
}

func (r *userCardSearchLogRepository) Update(auth *model.UserCardSearchLog) error {
	result := r.db.Save(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userCardSearchLogRepository) Delete(id string) error {
	result := r.db.Model(&model.UserCardSearchLog{}).Where("id = ?", id).Update("status", false)
	return result.Error
}

