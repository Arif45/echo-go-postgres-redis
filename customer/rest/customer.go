package rest

import (
	"fin-auth/domain"
	"fin-auth/dto"
	"fin-auth/models"
	"fin-auth/utils"
	"strings"

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

	if *req.VerificationType != utils.VERIFICATION_TYPE_RELIANCE && *req.VerificationType != utils.VERIFICATION_TYPE_STANDARD {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid verification type"))
	}

	var basicInfo = req.BasicInfo

	if !utils.IsValidName(basicInfo.FirstName) {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid first name"))
	}

	if !utils.IsValidName(basicInfo.LastName) {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid last name"))
	}

	if !utils.DateFieldValidation(basicInfo.DOB, "2006-01-02") {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid date of birth"))
	}

	if !utils.IsEmailValid(basicInfo.Email) {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid email address"))
	}

	_, err = utils.ValidateAndNormalizePhone(basicInfo.Phone, "")
	if err != nil {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid phone number"))
	}

	if strings.TrimSpace(basicInfo.CountryOfResidence) == " " {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid country of residence"))
	}

	if strings.TrimSpace(basicInfo.Nationality) == " " {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid nationality"))
	}
	if strings.TrimSpace(basicInfo.TIN) == " " {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid TIN"))
	}

	var address = req.Address

	if strings.TrimSpace(address.Street) == " " {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid street"))
	}
	if strings.TrimSpace(address.City) == " " {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid city"))
	}
	if strings.TrimSpace(address.State) == " " {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid state"))
	}
	if strings.TrimSpace(address.PostalCode) == " " {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid postal code"))
	}
	if strings.TrimSpace(address.Country) == " " {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid country"))
	}

	var financialProfile = req.FinancialProfile
	if financialProfile.OccupationID == nil || *financialProfile.OccupationID <= 0 {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid occupation ID"))
	}
	if financialProfile.SourceOfFundID == nil || *financialProfile.SourceOfFundID <= 0 {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid source of fund ID"))
	}
	if financialProfile.PurposeID == nil || *financialProfile.PurposeID <= 0 {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid purpose ID"))
	}
	if financialProfile.MonthlyVolumeUSD <= 0 {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid monthly volume in USD"))
	}
	if strings.TrimSpace(financialProfile.SOFDescription) == " " {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid source of fund description"))
	}

	var metaData = req.MetaData
	if strings.TrimSpace(metaData.Reference) == " " {
		return h.Response.InvalidData(c, utils.StringPtr("Invalid reference in meta data"))
	}
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
