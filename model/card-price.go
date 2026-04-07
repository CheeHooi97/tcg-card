package model

import (
	"pkm/utils"
	"time"
)

type CardPrice struct {
	Id       string `gorm:"primaryKey" json:"id"`
	CardId   string `json:"cardId"`
	Name     string `json:"name"`
	Set      string `json:"set"`
	Ungrade  string `json:"ungrade"`
	Grade7   string `json:"grade7"`
	Grade8   string `json:"grade8"`
	Grade9   string `json:"grade9"`
	Grade9_5 string `json:"grade9_5"`
	Grade10  string `json:"grade10"`
	BaseModel
}

func NewCardPrice() *CardPrice {
	now := time.Now().UTC()

	m := new(CardPrice)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *CardPrice) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *CardPrice) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
