package model

import "time"

type UserCardSearchLog struct {
	Id     string `gorm:"primaryKey" json:"id"`
	CardId string `json:"cardId"`
	BaseModel
}

func (m *UserCardSearchLog) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *UserCardSearchLog) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
