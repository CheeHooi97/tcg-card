package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	GetById(id string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userRepository) GetById(id string) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}

	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userRepository) Delete(id string) error {
	result := r.db.Model(&model.User{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
