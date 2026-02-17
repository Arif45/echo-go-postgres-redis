package models

import "time"

type RefreshToken struct {
	BaseModel
	ClientId      string    `gorm:"index" json:"client_id"`
	Token         string    `gorm:"uniqueIndex;size:100;not null" json:"token"`
	AccessTokenId uint      `json:"access_token_id"`
	ExpiredAt     time.Time `json:"expired_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
