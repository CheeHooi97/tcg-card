package model

import (
	"pkm/utils"
	"time"
)

type TAGUrl struct {
	Id  string `gorm:"primaryKey" json:"id"`
	Url string `json:"url"`
	BaseModel
}

func NewTAGUrl() *TAGUrl {
	now := time.Now().UTC()

	m := new(TAGUrl)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *TAGUrl) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *TAGUrl) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
