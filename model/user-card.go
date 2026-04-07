package model

import "time"

type UserCard struct {
	Id       string `gorm:"primaryKey" json:"id"`
	UserId   string `json:"userId"`
	CardId   string `json:"cardId"`
	AuthName string `json:"authName"`
	Quantity int64  `json:"quantity"`
	Grade    string `json:"grade"`
	BaseModel
}

func (m *UserCard) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *UserCard) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
