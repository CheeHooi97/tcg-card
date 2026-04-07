package handler

import (
	"net/http"
	"pkm/errcode"
	"pkm/service"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Admin             *service.AdminService
	BGS               *service.BGSService
	BGSUrl            *service.BGSUrlService
	BGSLogging        *service.BGSLoggingService
	Card              *service.CardService
	CardPrice         *service.CardPriceService
	CGC               *service.CGCService
	CGCUrl            *service.CGCUrlService
	CGCLogging        *service.CGCLoggingService
	PSA               *service.PSAService
	PSAUrl            *service.PSAUrlService
	PSALogging        *service.PSALoggingService
	Set               *service.SetService
	TAG               *service.TAGService
	TAGUrl            *service.TAGUrlService
	TAGLogging        *service.TAGLoggingService
	Token             *service.TokenService
	UserCardSearchLog *service.UserCardSearchLogService
	UserCard          *service.UserCardService
	User              *service.UserService
	UserDevice        *service.UserDeviceService
}

func NewHandler(services *service.Services) *Handler {
	h := &Handler{
		Admin:             services.AdminService,
		BGS:               services.BGSService,
		BGSUrl:            services.BGSUrlService,
		BGSLogging:        services.BGSLoggingService,
		Card:              services.CardService,
		CardPrice:         services.CardPriceService,
		CGC:               services.CGCService,
		CGCUrl:            services.CGCUrlService,
		CGCLogging:        services.CGCLoggingService,
		PSA:               services.PSAService,
		PSAUrl:            services.PSAUrlService,
		PSALogging:        services.PSALoggingService,
		Set:               services.SetService,
		TAG:               services.TAGService,
		TAGUrl:            services.TAGUrlService,
		TAGLogging:        services.TAGLoggingService,
		Token:             services.TokenService,
		UserCardSearchLog: services.UserCardSearchLogService,
		UserCard:          services.UserCardService,
		User:              services.UserService,
		UserDevice:        services.UserDeviceService,
	}

	return h
}

func responseError(c echo.Context, message errcode.ErrorCode) error {
	return c.JSON(http.StatusOK, map[string]any{
		"result": nil,
		"errmsg": message.Message,
		"error":  true,
		"status": false,
	})
}

func responseJSON(c echo.Context, result any) error {
	return c.JSON(http.StatusOK, map[string]any{
		"result": result,
		"errmsg": "",
		"error":  false,
		"status": true,
	})
}

func responseListJSON(c echo.Context, result any) error {
	return c.JSON(http.StatusOK, map[string]any{
		"result": map[string]any{
			"groups": result,
		},
		"errmsg": "",
		"error":  false,
		"status": true,
	})
}

func responseNoContent(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func responseValidationError(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, map[string]any{
		"result": nil,
		"errmsg": message,
		"error":  true,
		"status": false,
	})
}
