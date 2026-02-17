package models

type User struct {
	BaseModel
	ClientId    string `gorm:"uniqueIndex;size:100;not null" json:"client_id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Email       string `gorm:"uniqueIndex;not null" json:"email"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Description string `gorm:"size:100;default:'null'" json:"description"`
}

func (User) TableName() string {
	return "users"
}
