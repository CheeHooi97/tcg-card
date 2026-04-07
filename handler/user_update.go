package handler

import (
	"pkm/errcode"
	"pkm/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateUser(c echo.Context) error {
	id := c.Param("id")

	var i struct {
		CompanyId string `json:"companyId"`
		Username  string `json:"username"`
		FcmToken  string `json:"fcmToken"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	user, err := h.User.GetById(id)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	} else if user == nil {
		return responseError(c, errcode.FailedGetUser)
	}

	if i.CompanyId != "" {
		user.CompanyId = i.CompanyId
	}

	if i.Username != "" {
		user.Username = i.Username
	}

	if i.FcmToken != "" {
		user.FcmToken = i.FcmToken
	}

	user.UpdateDt()
	if err := h.User.Update(user); err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	return responseNoContent(c)
}
