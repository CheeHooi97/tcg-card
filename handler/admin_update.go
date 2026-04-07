package handler

import (
	"pkm/errcode"
	"pkm/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateAdmin(c echo.Context) error {
	id := c.Param("id")

	var i struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FcmToken string `json:"fcmToken"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	admin, err := h.Admin.GetAdminById(id)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	} else if admin == nil {
		return responseError(c, errcode.AdminNotFound)
	}

	if i.Username != "" {
		admin.Username = i.Username
	}

	if i.Email != "" {
		admin.Email = i.Email
	}

	if i.FcmToken != "" {
		admin.FcmToken = i.FcmToken
	}

	admin.UpdateDt()
	if err := h.Admin.UpdateAdmin(admin); err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	return responseNoContent(c)
}
