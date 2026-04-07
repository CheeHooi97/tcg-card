package handler

import (
	"pkm/errcode"
	"pkm/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetAdmin(c echo.Context) error {
	var i struct {
		Id string `json:"id" validate:"required"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	admin, err := h.Admin.GetAdminById(i.Id)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	} else if admin == nil {
		return responseError(c, errcode.AdminNotFound)
	}

	return responseJSON(c, admin)
}
