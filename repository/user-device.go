package repository

import (
	"errors"
	"pkm/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDeviceRepository interface {
	Create(device *model.UserDevice) error
	GetById(id string) (*model.UserDevice, error)
	GetAllByUserId(userId string) ([]*model.UserDevice, error)
	FindLastByUserId(userId string) (*model.UserDevice, error)
	FindByDeviceId(deviceId string) (*model.UserDevice, error)
	FindByUserIdAndDeviceID(userId, deviceId string) (*model.UserDevice, error)
	UpdateByPnsToken(token string) error
	Upsert(device *model.UserDevice) error
	Update(device *model.UserDevice) error
	Delete(id string) error
}

type userDeviceRepository struct {
	db *gorm.DB
}

func NewUserDeviceRepository(db *gorm.DB) UserDeviceRepository {
	return &userDeviceRepository{db: db}
}

func (r *userDeviceRepository) Create(device *model.UserDevice) error {
	result := r.db.Create(device)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userDeviceRepository) GetById(id string) (*model.UserDevice, error) {
	var device model.UserDevice
	result := r.db.First(&device, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &device, nil
}

func (r *userDeviceRepository) GetAllByUserId(userId string) ([]*model.UserDevice, error) {
	var devices []*model.UserDevice
	result := r.db.
		Where("user_id = ?", userId).
		Order("updated_date_time DESC").
		Find(&devices)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return devices, nil
}

func (r *userDeviceRepository) FindLastByUserId(userId string) (*model.UserDevice, error) {
	var device model.UserDevice
	result := r.db.
		Where("user_id = ?", userId).
		Order("updated_date_time DESC").
		First(&device)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &device, nil
}

func (r *userDeviceRepository) FindByDeviceId(deviceId string) (*model.UserDevice, error) {
	var device model.UserDevice
	result := r.db.
		Where("device_id = ?", deviceId).
		Order("created_date_time DESC").
		First(&device)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &device, nil
}

func (r *userDeviceRepository) FindByUserIdAndDeviceID(userId, deviceId string) (*model.UserDevice, error) {
	var device model.UserDevice
	result := r.db.
		Where("user_id = ? AND device_id != ?", userId, deviceId).
		Order("created_date_time DESC").
		First(&device)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &device, nil
}

func (r *userDeviceRepository) UpdateByPnsToken(token string) error {
	result := r.db.
		Model(&model.UserDevice{}).
		Where("pns_token = ?", token).
		Update("pns_token", "")

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
	}

	return nil
}

func (r *userDeviceRepository) Upsert(device *model.UserDevice) error {
	result := r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(device)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userDeviceRepository) Update(device *model.UserDevice) error {
	result := r.db.Save(device)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userDeviceRepository) Delete(id string) error {
	result := r.db.Model(&model.UserDevice{}).Where("id = ?", id).Updates(map[string]any{
		"status":          false,
		"updated_date_time": time.Now().UTC(),
	})
	return result.Error
}

