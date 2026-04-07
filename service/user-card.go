package service

import (
	"pkm/model"
	"pkm/repository"
)

type UserCardService struct {
	userCardRepo repository.UserCardRepository
}

func NewUserCardService(userCardRepo repository.UserCardRepository) *UserCardService {
	return &UserCardService{userCardRepo: userCardRepo}
}

func (s *UserCardService) Create(UserCard *model.UserCard) error {
	return s.userCardRepo.Create(UserCard)
}

func (s *UserCardService) GetById(id string) (*model.UserCard, error) {
	return s.userCardRepo.GetById(id)
}

func (s *UserCardService) GetAllCards() ([]*model.UserCard, error) {
	return s.userCardRepo.GetAllCards()
}

func (s *UserCardService) Update(UserCard *model.UserCard) error {
	return s.userCardRepo.Update(UserCard)
}

func (s *UserCardService) Delete(id string) error {
	return s.userCardRepo.Delete(id)
}
