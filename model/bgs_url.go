package model

import (
	"pkm/utils"
	"time"
)

type BGSUrl struct {
	Id  string `gorm:"primaryKey" json:"id"`
	Url string `json:"url"`
	BaseModel
}

func NewBGSUrl() *BGSUrl {
	now := time.Now().UTC()

	m := new(BGSUrl)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *BGSUrl) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *BGSUrl) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
