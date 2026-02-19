package rest

import (
	"fin-auth/domain"

	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	Service  domain.AddressService
	Response domain.Response
}

func SetupAddressRoutes(api *echo.Group, s domain.AddressService) {
	// handler := &AddressHandler{
	// 	Service:  s,
	// 	Response: domain.NewResponse(),
	// }
}
