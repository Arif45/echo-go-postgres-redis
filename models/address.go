package models

import "time"

type Address struct {
	ID         int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	CustomerID *string   `json:"customer_id" gorm:"column:customer_id;type:uuid;references:ID;constraint:OnDelete:CASCADE"`
	Customer   *Customer `json:"-" gorm:"foreignKey:CustomerID;references:ID"`
	Street     *string   `json:"street" gorm:"column:street"`
	City       *string   `json:"city" gorm:"column:city"`
	State      *string   `json:"state" gorm:"column:state"`
	StateCode  *string   `json:"state_code" gorm:"column:state_code"`
	PostalCode *string   `json:"postal_code" gorm:"column:postal_code"`
	Country    *string   `json:"country" gorm:"column:country"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (Address) TableName() string {
	return "addresses"
}
