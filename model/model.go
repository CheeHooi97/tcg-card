package model

import "time"

type BaseModel struct {
	CreatedDateTime time.Time `json:"createdAt"`
	UpdatedDateTime time.Time `json:"updatedAt"`
}
