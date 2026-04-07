package service

import (
	"pkm/model"
	"pkm/repository"
)

type UserDeviceService struct {
	userDeviceRepo repository.UserDeviceRepository
}

func NewUserDeviceService(userDeviceRepo repository.UserDeviceRepository) *UserDeviceService {
	return &UserDeviceService{userDeviceRepo: userDeviceRepo}
}

func (s *UserDeviceService) Create(device *model.UserDevice) error {
	return s.userDeviceRepo.Create(device)
}

func (s *UserDeviceService) GetById(id string) (*model.UserDevice, error) {
	return s.userDeviceRepo.GetById(id)
}

func (s *UserDeviceService) GetAllByUserId(userId string) ([]*model.UserDevice, error) {
	return s.userDeviceRepo.GetAllByUserId(userId)
}

func (s *UserDeviceService) FindLastByUserId(userId string) (*model.UserDevice, error) {
	return s.userDeviceRepo.FindLastByUserId(userId)
}

func (s *UserDeviceService) FindByDeviceId(deviceId string) (*model.UserDevice, error) {
	return s.userDeviceRepo.FindByDeviceId(deviceId)
}

func (s *UserDeviceService) FindByUserIdAndDeviceID(userId, deviceId string) (*model.UserDevice, error) {
	return s.userDeviceRepo.FindByUserIdAndDeviceID(userId, deviceId)
}

func (s *UserDeviceService) UpdateByPnsToken(token string) error {
	return s.userDeviceRepo.UpdateByPnsToken(token)
}

func (s *UserDeviceService) Upsert(device *model.UserDevice) error {
	return s.userDeviceRepo.Upsert(device)
}

func (s *UserDeviceService) Update(device *model.UserDevice) error {
	return s.userDeviceRepo.Update(device)
}

func (s *UserDeviceService) Delete(id string) error {
	return s.userDeviceRepo.Delete(id)
}
