package model

import (
	"pkm/utils"
	"time"
)

type PSAUrl struct {
	Id      string `gorm:"primaryKey" json:"id"`
	SetName string `json:"setName"`
	Url     string `json:"url"`
	BaseModel
}

func NewPSAUrl() *PSAUrl {
	now := time.Now().UTC()

	m := new(PSAUrl)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *PSAUrl) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *PSAUrl) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
