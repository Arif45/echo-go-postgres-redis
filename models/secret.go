package models

type Secret struct {
	BaseModel
	ClientId        string `json:"client_id"`
	Secret          string `gorm:"uniqueIndex;size:100;not null" json:"secret,omitempty"`
	SecondarySecret string `gorm:"uniqueIndex;size:100;null" json:"secondary_secret,omitempty"`
}

func (Secret) TableName() string {
	return "secrets"
}
