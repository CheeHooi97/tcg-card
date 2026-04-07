package handler

import (
	"pkm/errcode"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetAllAdmins(c echo.Context) error {
	admins, err := h.Admin.GetAllAdmins()
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	return responseJSON(c, admins)
}
