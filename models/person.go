package models

import "time"

type Person struct {
	ID                 int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	CustomerID         *string    `json:"customer_id" gorm:"column:customer_id;type:uuid;references:ID;constraint:OnDelete:CASCADE"`
	Customer           *Customer  `json:"-" gorm:"foreignKey:CustomerID;references:ID"`
	FirstName          *string    `json:"first_name" gorm:"column:first_name"`
	LastName           *string    `json:"last_name" gorm:"column:last_name"`
	Email              *string    `json:"email" gorm:"column:email"`
	Phone              *string    `json:"phone" gorm:"column:phone"`
	DOB                *time.Time `json:"dob" gorm:"column:dob"`
	CountryOfResidence *string    `json:"country_of_residence" gorm:"column:country_of_residence"`
	Nationality        *string    `json:"nationality" gorm:"column:nationality"`
	TIN                *string    `json:"tin" gorm:"column:tin"`
	Occupation         *string    `json:"occupation" gorm:"column:occupation"`
	SourceOfFundsID    *int       `json:"source_of_funds_id" gorm:"column:source_of_funds_id"`
	PurposeID          *int       `json:"purpose_id" gorm:"column:purpose_id"`
	Purpose            *string    `json:"purpose" gorm:"column:purpose"`
	MonthlyVolumeUSD   *float64   `json:"monthly_volume_usd" gorm:"column:monthly_volume_usd"`
	AddressID          *int       `json:"address_id" gorm:"column:address_id"`
	OccupationID       *int       `json:"occupation_id" gorm:"column:occupation_id"`
	SumsubApplicantID  *string    `json:"sumsub_applicant_id" gorm:"column:sumsub_applicant_id"`
	CreatedAt          time.Time  `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (Person) TableName() string {
	return "persons"
}
