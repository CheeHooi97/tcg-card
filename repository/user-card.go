package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type UserCardRepository interface {
	Create(auth *model.UserCard) error
	GetById(id string) (*model.UserCard, error)
	GetAllCards() ([]*model.UserCard, error)
	Update(auth *model.UserCard) error
	Delete(id string) error
}

type userCardRepository struct {
	db *gorm.DB
}

func NewUserCardRepository(db *gorm.DB) UserCardRepository {
	return &userCardRepository{db: db}
}

func (r *userCardRepository) Create(auth *model.UserCard) error {
	result := r.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userCardRepository) GetById(id string) (*model.UserCard, error) {
	var auth model.UserCard
	result := r.db.First(&auth, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &auth, nil
}

func (r *userCardRepository) GetAllCards() ([]*model.UserCard, error) {
	var UserCard []*model.UserCard
	result := r.db.
		Order("createdDateTime ASC").
		Find(&UserCard)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return UserCard, nil
}

func (r *userCardRepository) Update(auth *model.UserCard) error {
	result := r.db.Save(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userCardRepository) Delete(id string) error {
	result := r.db.Model(&model.UserCard{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
