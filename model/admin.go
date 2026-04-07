package model

import "time"

type Admin struct {
	Id        string `gorm:"primaryKey" json:"id"`
	CompanyId string `json:"companyId"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FcmToken  string `json:"fcmToken"`
	Status    bool   `json:"status"`
	BaseModel
}

func (m *Admin) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *Admin) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
