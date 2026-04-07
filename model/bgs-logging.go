package model

import (
	"pkm/utils"
	"time"
)

type BGSLogging struct {
	Id          string `gorm:"primaryKey" json:"id"`
	BGSId       string `json:"bgsId"`
	CardName    string `json:"cardName"`
	CardNumber  string `json:"cardNumber"`
	SetNumber   string `json:"setNumber"`
	SetName     string `json:"setName"`
	Rarity      string `json:"rarity"`
	Description string `json:"description"`
	SetId       string `json:"setId"`
	Total       string `json:"total"`
	Grade1      string `json:"Grade1"`
	Grade1_5    string `json:"Grade1_5"`
	Grade2      string `json:"Grade2"`
	Grade2_5    string `json:"Grade2_5"`
	Grade3      string `json:"Grade3"`
	Grade3_5    string `json:"Grade3_5"`
	Grade4      string `json:"Grade4"`
	Grade4_5    string `json:"Grade4_5"`
	Grade5      string `json:"Grade5"`
	Grade5_5    string `json:"Grade5_5"`
	Grade6      string `json:"Grade6"`
	Grade6_5    string `json:"Grade6_5"`
	Grade7      string `json:"Grade7"`
	Grade7_5    string `json:"Grade7_5"`
	Grade8      string `json:"Grade8"`
	Grade8_5    string `json:"Grade8_5"`
	Grade9      string `json:"Grade9"`
	Grade9_5    string `json:"Grade9_5"`
	Grade10P    string `json:"Grade10P"`
	Grade10BL   string `json:"Grade10BL"`
	BaseModel
}

func NewBGSLogging() *BGSLogging {
	now := time.Now().UTC()

	m := new(BGSLogging)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *BGSLogging) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *BGSLogging) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
