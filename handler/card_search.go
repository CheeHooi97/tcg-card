package handler

import (
	"pkm/errcode"
	"pkm/middleware"
	"pkm/model"
	"pkm/transformer"
	"pkm/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) SearchCard(c echo.Context) error {
	var i struct {
		Keyword string `json:"keyword"`
		Set     string `json:"set"`
		Order   string `json:"order"`
		UserId  string `json:"userId"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	actor, err := middleware.GetActor(c)
	if err != nil {
		return responseError(c, errcode.ActorNotFound)
	}

	user, err := h.User.GetById(actor.Id)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	} else if user.Id == "" {
		return responseError(c, errcode.UserNotFound)
	}

	var cards []*model.Card

	if i.Keyword != "" {
		cards, err = h.Card.SearchCardKeywords(i.Keyword)
		if err != nil {
			return responseError(c, errcode.InternalServerError)
		}
	} else {
		cards, err = h.Card.SearchCardBySort("relevance", user.Id)
		if err != nil {
			return responseError(c, errcode.InternalServerError)
		}
	}

	if i.Set != "" {

	}

	rarityMap := make(map[string]string, 0)
	increasedPriceMap := make(map[string]string)
	percentChangeMap := make(map[string]string)

	cardLists := transformer.ToCardLists(cards, rarityMap, increasedPriceMap, percentChangeMap)

	return responseJSON(c, echo.Map{
		"lists": cardLists,
	})
}
