package service

import (
	"pkm/model"
	"pkm/repository"
)

type AdminService struct {
	adminRepo repository.AdminRepository
}

func NewAdminService(adminRepo repository.AdminRepository) *AdminService {
	return &AdminService{adminRepo: adminRepo}
}

func (s *AdminService) CreateAdmin(admin *model.Admin) error {
	return s.adminRepo.CreateAdmin(admin)
}

func (s *AdminService) GetAdminById(id string) (*model.Admin, error) {
	return s.adminRepo.GetAdminById(id)
}

func (s *AdminService) GetAllAdmins() ([]*model.Admin, error) {
	return s.adminRepo.GetAllAdmins()
}

func (s *AdminService) UpdateAdmin(admin *model.Admin) error {
	return s.adminRepo.UpdateAdmin(admin)
}

func (s *AdminService) DeleteAdmin(id string) error {
	return s.adminRepo.DeleteAdmin(id)
}
