package model

import (
	"pkm/utils"
	"time"
)

type Card struct {
	Id        string `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	Number    string `json:"number"`
	SetNumber string `json:"setNumber"`
	SetName   string `json:"setName"`
	Rarity    string `json:"rarity"`
	Ungrade   string `json:"ungrade"`
	Grade7    string `json:"grade7"`
	Grade8    string `json:"grade8"`
	Grade9    string `json:"grade9"`
	Grade9_5  string `json:"grade9_5"`
	Grade10   string `json:"grade10"`
	PhotoUrl  string `json:"photoUrl"`
	BaseModel
}

func NewCard() *Card {
	now := time.Now().UTC()

	m := new(Card)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *Card) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *Card) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
