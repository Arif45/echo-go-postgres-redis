package rest

import (
	"fin-auth/domain"

	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	Service  domain.CustomerService
	Response domain.Response
}

func SetupCustomerRoutes(api *echo.Group, s domain.CustomerService) {
	// handler := &CustomerHandler{
	// 	Service:  s,
	// 	Response: domain.NewResponse(),
	// }

}
