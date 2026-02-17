package dto

import (
	"fin-auth/utils"
	"time"
)

type LoginReq struct {
	ClientId string `json:"client_id"`
	Secret   string `json:"secret"`
}

type TokenResponse struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

func (r *LoginReq) Validate() utils.Validation {
	v := utils.NewValidationError()
	errs := utils.ErrorResponse{}

	if r.ClientId == "" || !utils.StringFiledValidation(r.ClientId, 1, 100) {
		errs.Add("client_id", utils.ErrorMessage("client_id"))
		v.Status = true
	}

	if r.Secret == "" || !utils.StringFiledValidation(r.Secret, 1, 100) {
		errs.Add("secret", utils.ErrorMessage("secret"))
		v.Status = true
	}

	v.Response = errs
	return v
}
