package models

import "time"

type AccessToken struct {
	BaseModel
	ClientId  string    `gorm:"index" json:"client_id"`
	Token     string    `gorm:"uniqueIndex;size:100;not null" json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (AccessToken) TableName() string {
	return "access_tokens"
}
