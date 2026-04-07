package handler

import (
	"pkm/errcode"

	"github.com/labstack/echo/v4"
)

func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	if err := h.User.Delete(id); err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	return responseNoContent(c)
}
