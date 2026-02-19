package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Customer struct {
	ID                    string        `json:"id" gorm:"type:uuid;primaryKey"`
	ClientID              *string       `json:"client_id" gorm:"column:client_id"`
	CustomerType          *string       `json:"customer_type" gorm:"column:customer_type"`
	TOSPolicies           *string       `json:"tos_policies" gorm:"column:tos_policies;type:uuid"`
	USDEnable             *bool         `json:"usd_enable" gorm:"column:usd_enable"`
	KYCStatus             *string       `json:"kyc_status" gorm:"column:kyc_status"`
	VerificationType      string        `json:"verification_type" gorm:"column:verification_type"`
	BridgeKYCStatus       *string       `json:"bridge_kyc_status" gorm:"column:bridge_kyc_status"`
	CustomerStatus        *string       `json:"customer_status" gorm:"column:customer_status"`
	Meta                  *string       `json:"meta" gorm:"column:meta"`
	BridgeCustomerID      *string       `json:"bridge_customer_id" gorm:"column:bridge_customer_id"`
	AvailableCorridorsIDs pq.Int64Array `json:"available_corridors_ids" gorm:"column:available_corridors_ids;type:integer[]"`
	CreatedAt             time.Time     `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt             time.Time     `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (Customer) TableName() string {
	return "customers"
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
