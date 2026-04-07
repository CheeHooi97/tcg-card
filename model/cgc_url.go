package model

import (
	"pkm/utils"
	"time"
)

type CGCUrl struct {
	Id  string `gorm:"primaryKey" json:"id"`
	Url string `json:"url"`
	BaseModel
}

func NewCGCUrl() *CGCUrl {
	now := time.Now().UTC()

	m := new(CGCUrl)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *CGCUrl) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *CGCUrl) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
