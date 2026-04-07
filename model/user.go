package model

import (
	"pkm/utils"
	"time"
)

type User struct {
	Id        string `gorm:"primaryKey" json:"id"`
	CompanyId string `json:"companyId"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	PhotoURL  string `json:"photoUrl" sqlike:",longtext"`
	FcmToken  string `json:"fcmToken"`
	Status    bool   `json:"status"`
	BaseModel
}

func NewUser() *User {
	now := time.Now().UTC()

	m := new(User)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *User) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *User) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
