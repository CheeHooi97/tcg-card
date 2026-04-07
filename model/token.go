package model

import (
	"pkm/utils"
	"time"
)

type Token struct {
	Id          string `gorm:"primaryKey" json:"id"`
	ReferenceId string `json:"referenceId"`
	DeviceId    string `json:"deviceId"`
	AccessToken string `json:"accessToken" gorm:"type:LONGTEXT"`
	BaseModel
}

func NewToken() *Token {
	now := time.Now().UTC()

	m := new(Token)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *Token) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}

type UserWithToken struct {
	*User
	Token string `json:"token"`
}
