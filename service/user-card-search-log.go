package service

import (
	"pkm/model"
	"pkm/repository"
)

type UserCardSearchLogService struct {
	userCardSearchLogRepo repository.UserCardSearchLogRepository
}

func NewUserCardSearchLogService(userCardSearchLogRepo repository.UserCardSearchLogRepository) *UserCardSearchLogService {
	return &UserCardSearchLogService{userCardSearchLogRepo: userCardSearchLogRepo}
}

func (s *UserCardSearchLogService) Create(UserCardSearchLog *model.UserCardSearchLog) error {
	return s.userCardSearchLogRepo.Create(UserCardSearchLog)
}

func (s *UserCardSearchLogService) GetById(id string) (*model.UserCardSearchLog, error) {
	return s.userCardSearchLogRepo.GetById(id)
}

func (s *UserCardSearchLogService) GetAllCards() ([]*model.UserCardSearchLog, error) {
	return s.userCardSearchLogRepo.GetAllCards()
}

func (s *UserCardSearchLogService) Update(UserCardSearchLog *model.UserCardSearchLog) error {
	return s.userCardSearchLogRepo.Update(UserCardSearchLog)
}

func (s *UserCardSearchLogService) Delete(id string) error {
	return s.userCardSearchLogRepo.Delete(id)
}
