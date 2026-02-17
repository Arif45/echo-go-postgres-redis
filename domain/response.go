package domain

import (
	"fin-auth/utils"

	"github.com/labstack/echo/v4"
)

type Response interface {
	Ok(c echo.Context) error
	SuccessOk(echo.Context, interface{}) error
	SuccessMessage(echo.Context, string) error

	// Validation block
	UnauthorizedResponse(c echo.Context, debug *string) error
	InternalServerError(c echo.Context, err error) error
	ConflictError(c echo.Context, message *string, err map[string]interface{}) error
	ValidationFail(c echo.Context, err utils.Validation, message *string) error
	CustomValidationFail(c echo.Context, err interface{}, message string, statusCode ...int) error
	InvalidData(c echo.Context, message *string) error
	ForbiddenResponse(c echo.Context, err error, message *string) error
	LockedResponse(c echo.Context, message *string) error
	NotFound(c echo.Context, message *string) error
}
