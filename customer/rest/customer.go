package rest

import (
	"fin-auth/domain"
	"fin-auth/dto"
	"fin-auth/models"
	"fin-auth/utils"

	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	Service  domain.CustomerService
	Response domain.Response
}

func SetupCustomerRoutes(api *echo.Group, s domain.CustomerService) {
	handler := &CustomerHandler{
		Service:  s,
		Response: domain.NewResponse(),
	}
	customer := api.Group("/customers")
	customer.POST("/individual", handler.createIndividualCustomer)
	customer.GET("", handler.listCustomers)
}

func (h *CustomerHandler) createIndividualCustomer(c echo.Context) error {
	var req dto.CreateCustomerRequest
	err := c.Bind(&req)
	if err != nil {
		return h.Response.InvalidData(c, nil)
	}

	errMessage := req.ValidateRequest()
	if errMessage != nil {
		return h.Response.InvalidData(c, errMessage)
	}
	var basicInfo = req.BasicInfo
	var address = req.Address
	var customer = req.GetCustomer(&basicInfo, utils.StringPtr(c.Get("client_id").(string)))
	var person = req.GetPerson(&basicInfo)
	var addressModel *models.Address
	addressModel = req.GetAddress(&address)

	res, _, _, err := h.Service.CreateIndividualCustomer(c.Request().Context(), customer, person, addressModel, &req)

	if err != nil {
		return h.Response.InternalServerError(c, err)
	}
	return h.Response.SuccessOk(c, res)
}

func (h *CustomerHandler) listCustomers(c echo.Context) error {
	clientId := c.Get("client_id").(string)
	customers, err := h.Service.GetCustomersByClientID(c.Request().Context(), clientId)
	if err != nil {
		return h.Response.InternalServerError(c, err)
	}
	return h.Response.SuccessOk(c, customers)
}
