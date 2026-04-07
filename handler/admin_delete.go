package handler

import (
	"pkm/errcode"

	"github.com/labstack/echo/v4"
)

func (h *Handler) DeleteAdmin(c echo.Context) error {
	id := c.Param("id")

	if err := h.Admin.DeleteAdmin(id); err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	return responseNoContent(c)
}
