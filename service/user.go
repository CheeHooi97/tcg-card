package service

import (
	"pkm/model"
	"pkm/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Create(user *model.User) error {
	return s.userRepo.Create(user)
}

func (s *UserService) GetById(id string) (*model.User, error) {
	return s.userRepo.GetById(id)
}

func (s *UserService) GetByEmail(email string) (*model.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *UserService) Update(user *model.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) Delete(id string) error {
	return s.userRepo.Delete(id)
}
