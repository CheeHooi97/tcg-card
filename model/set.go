package model

import (
	"pkm/utils"
	"time"
)

type Set struct {
	Id       string `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	PhotoUrl string `json:"photoUrl"`
	BaseModel
}

func NewSet() *Set {
	now := time.Now().UTC()

	m := new(Set)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *Set) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *Set) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
