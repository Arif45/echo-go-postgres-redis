package domain

import (
	"fin-auth/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseImpl struct{}

func NewResponse() Response {
	return &ResponseImpl{}
}

func (r *ResponseImpl) Ok(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func (r *ResponseImpl) SuccessOk(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func (r *ResponseImpl) SuccessMessage(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": message,
	})
}

func (r *ResponseImpl) UnauthorizedResponse(c echo.Context, debug *string) error {
	response := map[string]interface{}{
		"success": false,
		"message": "Unauthorized access",
	}
	if debug != nil {
		response["debug"] = *debug
	}
	return c.JSON(http.StatusUnauthorized, response)
}

func (r *ResponseImpl) InternalServerError(c echo.Context, err error) error {
	response := map[string]interface{}{
		"success": false,
		"message": "Internal server error",
	}
	if err != nil {
		response["error"] = err.Error()
	}
	return c.JSON(http.StatusInternalServerError, response)
}

func (r *ResponseImpl) ConflictError(c echo.Context, message *string, err map[string]interface{}) error {
	msg := "Conflict error"
	if message != nil {
		msg = *message
	}
	response := map[string]interface{}{
		"success": false,
		"message": msg,
	}
	if err != nil {
		response["errors"] = err
	}
	return c.JSON(http.StatusConflict, response)
}

func (r *ResponseImpl) ValidationFail(c echo.Context, err utils.Validation, message *string) error {
	msg := err.Message
	if message != nil {
		msg = *message
	}
	return c.JSON(err.StatusCode, map[string]interface{}{
		"success": false,
		"message": msg,
		"errors":  err.Response,
	})
}

func (r *ResponseImpl) CustomValidationFail(c echo.Context, err interface{}, message string, statusCode ...int) error {
	code := http.StatusUnprocessableEntity
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	return c.JSON(code, map[string]interface{}{
		"success": false,
		"message": message,
		"errors":  err,
	})
}

func (r *ResponseImpl) InvalidData(c echo.Context, message *string) error {
	msg := "Invalid data provided"
	if message != nil {
		msg = *message
	}
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"success": false,
		"message": msg,
	})
}

func (r *ResponseImpl) ForbiddenResponse(c echo.Context, err error, message *string) error {
	msg := "Forbidden"
	if message != nil {
		msg = *message
	}
	response := map[string]interface{}{
		"success": false,
		"message": msg,
	}
	if err != nil {
		response["error"] = err.Error()
	}
	return c.JSON(http.StatusForbidden, response)
}

func (r *ResponseImpl) LockedResponse(c echo.Context, message *string) error {
	msg := "Resource locked"
	if message != nil {
		msg = *message
	}
	return c.JSON(http.StatusLocked, map[string]interface{}{
		"success": false,
		"message": msg,
	})
}

func (r *ResponseImpl) NotFound(c echo.Context, message *string) error {
	msg := "Resource not found"
	if message != nil {
		msg = *message
	}
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"success": false,
		"message": msg,
	})
}
