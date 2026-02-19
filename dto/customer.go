package dto

import (
	"fin-auth/models"
	"fin-auth/utils"
	"strings"
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

func (req *CreateCustomerRequest) ValidateRequest() *string {
	if *req.VerificationType != utils.VERIFICATION_TYPE_RELIANCE && *req.VerificationType != utils.VERIFICATION_TYPE_STANDARD {
		return utils.StringPtr("Invalid verification type")
	}

	var basicInfo = req.BasicInfo

	if !utils.IsValidName(basicInfo.FirstName) {
		return utils.StringPtr("Invalid first name")
	}

	if !utils.IsValidName(basicInfo.LastName) {
		return utils.StringPtr("Invalid last name")
	}

	if !utils.DateFieldValidation(basicInfo.DOB, "2006-01-02") {
		return utils.StringPtr("Invalid date of birth")
	}

	if !utils.IsEmailValid(basicInfo.Email) {
		return utils.StringPtr("Invalid email address")
	}

	_, err := utils.ValidateAndNormalizePhone(basicInfo.Phone, "")
	if err != nil {
		return utils.StringPtr("Invalid phone number")
	}

	if strings.TrimSpace(basicInfo.CountryOfResidence) == " " {
		return utils.StringPtr("Invalid country of residence")
	}

	if strings.TrimSpace(basicInfo.Nationality) == " " {
		return utils.StringPtr("Invalid nationality")
	}
	if strings.TrimSpace(basicInfo.TIN) == " " {
		return utils.StringPtr("Invalid TIN")
	}

	var address = req.Address

	if strings.TrimSpace(address.Street) == " " {
		return utils.StringPtr("Invalid street")
	}
	if strings.TrimSpace(address.City) == " " {
		return utils.StringPtr("Invalid city")
	}
	if strings.TrimSpace(address.State) == " " {
		return utils.StringPtr("Invalid state")
	}
	if strings.TrimSpace(address.PostalCode) == " " {
		return utils.StringPtr("Invalid postal code")
	}
	if strings.TrimSpace(address.Country) == " " {
		return utils.StringPtr("Invalid country")
	}

	var financialProfile = req.FinancialProfile
	if financialProfile.OccupationID == nil || *financialProfile.OccupationID <= 0 {
		return utils.StringPtr("Invalid occupation ID")
	}
	if financialProfile.SourceOfFundID == nil || *financialProfile.SourceOfFundID <= 0 {
		return utils.StringPtr("Invalid source of fund ID")
	}
	if financialProfile.PurposeID == nil || *financialProfile.PurposeID <= 0 {
		return utils.StringPtr("Invalid purpose ID")
	}
	if financialProfile.MonthlyVolumeUSD <= 0 {
		return utils.StringPtr("Invalid monthly volume in USD")
	}
	if strings.TrimSpace(financialProfile.SOFDescription) == " " {
		return utils.StringPtr("Invalid source of fund description")
	}

	var metaData = req.MetaData
	if strings.TrimSpace(metaData.Reference) == " " {
		return utils.StringPtr("Invalid reference in meta data")
	}
	return nil
}
