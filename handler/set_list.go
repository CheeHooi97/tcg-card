package handler

import (
	"pkm/config"
	"pkm/errcode"
	"pkm/kit/oss"
	"pkm/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) SetList(c echo.Context) error {
	var i struct {
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	lists, _ := h.Set.GetById("1769071018360741846")

	url, err := oss.RetrieveSignedURL(config.OSSBucket, lists.PhotoUrl)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	lists.PhotoUrl = url

	return responseJSON(c, echo.Map{
		"lists": lists,
	})
}
