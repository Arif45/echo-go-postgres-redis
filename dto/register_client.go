package dto

import (
	"fin-auth/utils"
)

type RegisterClientReq struct {
	Name        string `gorm:"size:100;not null" json:"name"`
	Email       string `gorm:"uniqueIndex;not null" json:"email"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Description string `gorm:"size:100;default:'null'" json:"description"`
}

type RegisterClientRes struct {
	ClientId        string `json:"client_id"`
	Secret          string `gorm:"size:100;not null" json:"secret,omitempty"`
	SecondarySecret string `gorm:"size:100;null" json:"secondary_secret,omitempty"`
}

func (r *RegisterClientReq) Validate() utils.Validation {
	v := utils.NewValidationError()
	errs := utils.ErrorResponse{}

	if r.Name == "" || !utils.StringFiledValidation(r.Name, 1, 50) {
		errs.Add("name", utils.ErrorMessage("name"))
		v.Status = true
	}

	if r.Email == "" || !utils.IsEmailValid(r.Email) {
		errs.Add("email", utils.ErrorMessage("email"))
		v.Status = true
	}

	v.Response = errs
	return v
}
