package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	ErrNotFound                = errors.New("resource not found")
	ErrInvalidPage             = errors.New("invalid page request")
	ErrConflict                = errors.New("data conflict or already exist")
	ErrBadRequest              = errors.New("bad request, check param or body")
	ErrForbidden               = errors.New("forbidden")
	ErrUnprocessableEntity     = errors.New("action could not be processed properly due to invalid data provided")
	ErrUnauthorized            = errors.New("unauthorized: you do not have permission to access this resource")
	ErrUnauthenticated         = errors.New("unauthenticated: no authenticated user found")
	ErrCountryMismatchPOA      = errors.New("country of residence does not match with proof of address")
	ErrCustomerNotFound        = errors.New("appropriate customer_id required")
	NoOrganizationFound        = errors.New("No organization found for this user")
	ErrFxRateNotFound          = errors.New("fx rate not found for the given currency pair")
	ErrFeeCalcMaxAmount        = errors.New("maximum amount exceeded for fee calculation")
	ErrFeeExceedsAmount        = errors.New("calculated fee exceeds the amount")
	ErrSrcAmountBelowMin       = errors.New("source amount is below the minimum required amount")
	ErrSrcAmountZeroOrNeg      = errors.New("source amount must be greater than zero")
	ErrLaDrainNotEnoughBalance = errors.New("Insufficient balance on liquidation address.")

	// LaDrainNotEnoughBalance = "Insufficient {source.currency} balance on {source.rail} liquidation address to settle {source_amount} {source.currency}. Please fund the liquidation address"
	LaDrainNotEnoughBalance   = "Insufficient %s balance on %s liquidation address to settle %s %s. Please fund the liquidation address."
	PreFundedNotEnoughBalance = "Pre-funded wallet on %s does not have enough balance to settle %s %s. Please add funds."

	InvalidFmt       = "invalid %s"
	InvalidChoiceFmt = "%s must be among %v"
	NoGroupFmt       = "no %s group found for promo %s"

	// Usage
	HourLockExpiredError = "this promo is only valid from %s to %s"
	MinOrderVal          = "please order %s %.0f (after discount) or more to apply this promo"
	Xtrip                = "Sorry, this promo is applicable for first %d orders only"
	MaxAttempt           = "Max otp attempt"
	MaxGetOtp            = "Max get otp"
)

func GetStatusCode(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInvalidPage:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrUnprocessableEntity:
		return http.StatusUnprocessableEntity
	case ErrUnauthorized:
		return http.StatusForbidden
	case ErrUnauthenticated:
		return http.StatusUnauthorized
	default:
		wrapErr := &WrapErr{}
		if errors.As(err, wrapErr) {
			return wrapErr.StatusCode
		}
		return http.StatusInternalServerError
	}
}

type InvalidTransitionError struct {
	FromStatus string
	FromSub    string
	ToStatus   string
	ToSub      string
}

func (e *InvalidTransitionError) Error() string {
	return fmt.Sprintf("invalid status transition: %s/%s → %s/%s",
		e.FromStatus, e.FromSub, e.ToStatus, e.ToSub)
}

func GetErrCode(err error) string {
	wrapErr := &WrapErr{}
	if errors.As(err, wrapErr) {
		return wrapErr.ErrCode
	}
	return ""
}

type WrapErr struct {
	Err        error
	StatusCode int
	ErrCode    string
}

func (e WrapErr) Error() string {
	return e.Err.Error()
}

func (e WrapErr) Unwrap() error {
	return e.Err // Returns inner error
}

func WrapError(err error, statusCode int, errCode string) error {
	return WrapErr{
		Err:        err,
		ErrCode:    errCode,
		StatusCode: statusCode,
	}
}

func IsDuplicateKeyError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// Check PostgreSQL specific errors
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return true
	}

	// String-based fallback
	errStr := err.Error()
	return strings.Contains(errStr, "duplicate key value violates unique constraint") ||
		strings.Contains(errStr, "SQLSTATE 23505")
}
