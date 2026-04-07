package handler

import (
	"pkm/errcode"
	"pkm/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUser(c echo.Context) error {
	var i struct {
		Id string `json:"id" validate:"required"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	user, err := h.User.GetById(i.Id)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	} else if user == nil {
		return responseError(c, errcode.UserNotFound)
	}

	return responseJSON(c, user)
}
