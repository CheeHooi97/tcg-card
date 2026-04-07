package repository

import (
	"errors"
	"pkm/model"

	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(admin *model.Admin) error
	GetAdminById(id string) (*model.Admin, error)
	GetAllAdmins() ([]*model.Admin, error)
	UpdateAdmin(admin *model.Admin) error
	DeleteAdmin(id string) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) CreateAdmin(admin *model.Admin) error {
	result := r.db.Create(admin)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *adminRepository) GetAdminById(id string) (*model.Admin, error) {
	var admin model.Admin
	result := r.db.First(&admin, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &admin, nil
}

func (r *adminRepository) GetAllAdmins() ([]*model.Admin, error) {
	var admins []*model.Admin
	result := r.db.Find(&admins)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return admins, nil
}

func (r *adminRepository) UpdateAdmin(admin *model.Admin) error {
	result := r.db.Save(admin)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *adminRepository) DeleteAdmin(id string) error {
	result := r.db.Model(&model.Admin{}).Where("id = ?", id).Update("status", false)
	return result.Error
}
