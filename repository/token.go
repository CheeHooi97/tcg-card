package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type TokenRepository interface {
	Create(token *model.Token) error
	GetById(id string) (*model.Token, error)
	FindByReferenceIdAndDeviceId(referenceId, deviceId string) (*model.Token, error)
	Update(token *model.Token) error
	Delete(id string) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(token *model.Token) error {
	result := r.db.Create(token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *tokenRepository) GetById(id string) (*model.Token, error) {
	var token model.Token
	result := r.db.First(&token, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &token, nil
}

func (r *tokenRepository) FindByReferenceIdAndDeviceId(referenceId, deviceId string) (*model.Token, error) {
	var token model.Token
	result := r.db.
		Where("reference_id = ? AND device_id = ?", referenceId, deviceId).
		First(&token)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &token, nil
}

func (r *tokenRepository) Update(token *model.Token) error {
	result := r.db.Save(token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *tokenRepository) Delete(id string) error {
	result := r.db.Model(&model.Token{}).Where("id = ?", id).Update("status", false)
	return result.Error
}

