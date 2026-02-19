package dto

import (
	"fin-auth/models"
	"time"

	"github.com/google/uuid"
)

type CreateCustomerRequest struct {
	VerificationType *string          `json:"verification_type"`
	BasicInfo        BasicInfo        `json:"basic_info" validate:"required"`
	Address          AddressInfo      `json:"address" validate:"required"`
	FinancialProfile FinancialProfile `json:"financial_profile" validate:"required"`
	MetaData         MetaData         `json:"meta_data" validate:"required"`
}

type BasicInfo struct {
	FirstName          string `json:"first_name" validate:"required"`
	LastName           string `json:"last_name" validate:"required"`
	DOB                string `json:"dob" validate:"required"`
	Email              string `json:"email" validate:"required"`
	Phone              string `json:"phone" validate:"required"`
	CountryOfResidence string `json:"country_of_residence" validate:"required"`
	Nationality        string `json:"nationality" validate:"required"`
	TIN                string `json:"tin" validate:"required"`
}

type AddressInfo struct {
	Street     string `json:"street" validate:"required"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
	Country    string `json:"country" validate:"required"`
}

type FinancialProfile struct {
	OccupationID     *int    `json:"occupation_id" validate:"required"`
	SourceOfFundID   *int    `json:"source_of_fund_id" validate:"required"`
	PurposeID        *int    `json:"purpose_id" validate:"required"`
	MonthlyVolumeUSD float64 `json:"monthly_volume_usd" validate:"required"`
	SOFDescription   string  `json:"sof_description" validate:"required"`
}

type MetaData struct {
	Reference string `json:"reference" validate:"required"`
}

func (r *CreateCustomerRequest) GetCustomer(basicInfo *BasicInfo, clientId *string) *models.Customer {
	return &models.Customer{
		ClientID:         clientId,
		CustomerType:     ptrString("individual"),
		TOSPolicies:      ptrString(uuid.New().String()),
		USDEnable:        nil,
		KYCStatus:        ptrString("INCOMPLETE"),
		VerificationType: *r.VerificationType,
		BridgeKYCStatus:  nil,
		CustomerStatus:   nil,
		Meta:             ptrString(r.MetaData.Reference),
	}
}

func (r *CreateCustomerRequest) GetPerson(basicInfo *BasicInfo) *models.Person {
	var dob *time.Time
	if basicInfo.DOB != "" {
		parsedDOB, err := time.Parse("2006-01-02", basicInfo.DOB)
		if err == nil {
			dob = &parsedDOB
		}
	}
	return &models.Person{
		CustomerID:         ptrString(uuid.New().String()),
		FirstName:          ptrString(basicInfo.FirstName),
		LastName:           ptrString(basicInfo.LastName),
		DOB:                dob,
		Email:              ptrString(basicInfo.Email),
		Phone:              ptrString(basicInfo.Phone),
		CountryOfResidence: ptrString(basicInfo.CountryOfResidence),
		Nationality:        ptrString(basicInfo.Nationality),
		TIN:                ptrString(basicInfo.TIN),
	}
}

func (r *CreateCustomerRequest) GetAddress(addressInfo *AddressInfo) *models.Address {
	return &models.Address{
		Street:     ptrString(addressInfo.Street),
		City:       ptrString(addressInfo.City),
		State:      ptrString(addressInfo.State),
		PostalCode: ptrString(addressInfo.PostalCode),
		Country:    ptrString(addressInfo.Country),
	}
}

// ptrString returns a pointer to the given string.
func ptrString(s string) *string {
	return &s
}
