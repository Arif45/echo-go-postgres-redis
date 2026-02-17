package dto

import (
	"fin-auth/utils"
	"time"
)

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRes struct {
	AccessToken     string    `json:"access_token"`
	AccessExpiresAt time.Time `json:"access_expires_at"`
}

func (r *RefreshTokenReq) Validate() utils.Validation {
	v := utils.NewValidationError()
	errs := utils.ErrorResponse{}

	if r.RefreshToken == "" || !utils.StringFiledValidation(r.RefreshToken, 1, 100) {
		errs.Add("refresh_token", utils.ErrorMessage("refresh_token"))
		v.Status = true
	}

	v.Response = errs
	return v
}
