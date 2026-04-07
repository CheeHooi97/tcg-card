package model

import (
	"pkm/utils"
	"time"
)

type PSA struct {
	Id          string `gorm:"primaryKey" json:"id"`
	CardName    string `json:"cardName"`
	CardNumber  string `json:"cardNumber"`
	SetNumber   string `json:"setNumber"`
	SetName     string `json:"setName"`
	Rarity      string `json:"rarity"`
	Description string `json:"description"`
	SpecID      string `json:"specId"`
	Total       string `json:"total"`
	Auth        string `json:"Auth"`
	Grade1      string `json:"Grade1"`
	Grade2      string `json:"Grade2"`
	Grade3      string `json:"Grade3"`
	Grade4      string `json:"Grade4"`
	Grade5      string `json:"Grade5"`
	Grade6      string `json:"Grade6"`
	Grade7      string `json:"Grade7"`
	Grade8      string `json:"Grade8"`
	Grade9      string `json:"Grade9"`
	Grade10     string `json:"Grade10"`
	BaseModel
}

func NewPSA() *PSA {
	now := time.Now().UTC()

	m := new(PSA)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *PSA) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *PSA) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
