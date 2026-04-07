package service

import (
	"pkm/model"
	"pkm/repository"
)

type TokenService struct {
	TokenRepo repository.TokenRepository
}

func NewTokenService(TokenRepo repository.TokenRepository) *TokenService {
	return &TokenService{TokenRepo: TokenRepo}
}

func (s *TokenService) Create(Token *model.Token) error {
	return s.TokenRepo.Create(Token)
}

func (s *TokenService) GetById(id string) (*model.Token, error) {
	return s.TokenRepo.GetById(id)
}

func (s *TokenService) FindByReferenceIdAndDeviceId(referenceId, deviceId string) (*model.Token, error) {
	return s.TokenRepo.FindByReferenceIdAndDeviceId(referenceId, deviceId)
}

func (s *TokenService) Update(Token *model.Token) error {
	return s.TokenRepo.Update(Token)
}

func (s *TokenService) Delete(id string) error {
	return s.TokenRepo.Delete(id)
}
