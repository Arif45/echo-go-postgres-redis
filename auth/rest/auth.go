package rest

import (
	"fin-auth/domain"
	"fin-auth/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Service  domain.AuthService
	Response domain.Response
}

func SetupAuthRoutes(api *echo.Group, s domain.AuthService) {
	handler := &AuthHandler{
		Service:  s,
		Response: domain.NewResponse(),
	}

	auth := api.Group("/auth")

	auth.POST("/register", handler.register)
	auth.POST("/login", handler.login)
	auth.POST("/refresh", handler.refreshToken)
}

func SetupProtectedRoutes(api *echo.Group, s domain.AuthService) {
	handler := &AuthHandler{
		Service:  s,
		Response: domain.NewResponse(),
	}

	api.GET("/auth/me", handler.me)
	api.POST("/auth/logout", handler.logout)
	api.GET("/auth/sessions", handler.listSessions)
	api.DELETE("/auth/sessions/:token", handler.revokeSession)
}

func (authHandler *AuthHandler) register(c echo.Context) error {
	req := dto.RegisterClientReq{}
	err := c.Bind(&req)

	if err != nil {
		return authHandler.Response.InvalidData(c, nil)
	}

	v := req.Validate()
	if v.Status {
		return authHandler.Response.ValidationFail(c, v, nil)
	}

	res, err := authHandler.Service.RegisterClient(c.Request().Context(), &req)

	if err != nil {
		return authHandler.Response.InternalServerError(c, err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Client registered successfully",
		"data":    res,
	})
}

func (authHandler *AuthHandler) login(c echo.Context) error {
	req := dto.LoginReq{}
	err := c.Bind(&req)
	if err != nil {
		return authHandler.Response.InvalidData(c, nil)
	}

	v := req.Validate()
	if v.Status {
		return authHandler.Response.ValidationFail(c, v, nil)
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	clientDatares, res, err := authHandler.Service.Login(c.Request().Context(), &req, ipAddress, userAgent)
	if clientDatares != nil {
		if req.ClientId != clientDatares.Secret.ClientId {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"message": "Invalid client ID",
			})
		}
		if req.Secret != clientDatares.Secret.Secret && req.Secret != clientDatares.Secret.SecondarySecret {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"message": "Invalid secret",
			})
		}

	}
	if err != nil {
		return authHandler.Response.InternalServerError(c, err)
	}

	return authHandler.Response.SuccessOk(c, res)
}

func (authHandler *AuthHandler) refreshToken(c echo.Context) error {
	req := dto.RefreshTokenReq{}
	err := c.Bind(&req)
	if err != nil {
		return authHandler.Response.InvalidData(c, nil)
	}
	v := req.Validate()
	if v.Status {
		return authHandler.Response.ValidationFail(c, v, nil)
	}

	res, err := authHandler.Service.RefreshToken(c.Request().Context(), &req)
	if err != nil {
		return authHandler.Response.InternalServerError(c, err)
	}

	return authHandler.Response.SuccessOk(c, res)
}

func (authHandler *AuthHandler) me(c echo.Context) error {
	clientId := c.Get("client_id")
	if clientId == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":   true,
		"client_id": clientId,
		"message":   "Authenticated client",
	})
}

func (authHandler *AuthHandler) logout(c echo.Context) error {
	token := c.Get("token")
	if token == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "No token found",
		})
	}

	err := authHandler.Service.Logout(c.Request().Context(), token.(string))
	if err != nil {
		return authHandler.Response.InternalServerError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Successfully logged out",
	})
}

func (authHandler *AuthHandler) listSessions(c echo.Context) error {
	clientId := c.Get("client_id")
	if clientId == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	sessions, err := authHandler.Service.ListSessions(c.Request().Context(), clientId.(string))
	if err != nil {
		return authHandler.Response.InternalServerError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":  true,
		"sessions": sessions,
		"count":    len(sessions),
	})
}

func (authHandler *AuthHandler) revokeSession(c echo.Context) error {
	clientId := c.Get("client_id")
	if clientId == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	tokenToRevoke := c.Param("token")
	if tokenToRevoke == "" {
		return authHandler.Response.InvalidData(c, nil)
	}

	err := authHandler.Service.RevokeSession(c.Request().Context(), clientId.(string), tokenToRevoke)
	if err != nil {
		return authHandler.Response.InternalServerError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Session revoked successfully",
	})
}
