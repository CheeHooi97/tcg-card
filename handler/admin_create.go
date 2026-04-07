package handler

import (
	"pkm/errcode"
	"pkm/model"
	"pkm/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateAdmin(c echo.Context) error {
	var i struct {
		Username  string `json:"username" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		CompanyId string `json:"companyId"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	admin := new(model.Admin)
	admin.Id = utils.UniqueID()
	admin.Username = i.Username
	admin.Email = i.Email
	admin.CompanyId = i.CompanyId
	admin.Status = true
	admin.DateTime()

	if err := h.Admin.CreateAdmin(admin); err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	return responseJSON(c, admin)
}
