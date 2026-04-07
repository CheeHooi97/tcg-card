package handler

import (
	"fmt"
	"pkm/errcode"
	"pkm/transformer"
	"pkm/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CardDetail(c echo.Context) error {
	var i struct {
		Id string `json:"id" validate:"required"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	// actor, err := middleware.GetActor(c)
	// if err != nil {
	// 	return responseError(c, errcode.ActorNotFound)
	// }

	// user, err := h.User.GetById(actor.Id)
	// if err != nil {
	// 	return responseError(c, errcode.InternalServerError)
	// } else if user.Id == "" {
	// 	return responseError(c, errcode.UserNotFound)
	// }

	card, err := h.Card.GetById(i.Id)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}
	var rarity string
	var increasedPrice string
	var percentChange string

	// cardName := "Swinub #31"
	var name, number string
	parts := strings.Split(card.Name, "#")
	if len(parts) == 2 {
		name = strings.TrimSpace(parts[0])
		number = strings.TrimSpace(parts[1])
	}

	card.Name = name
	card.Number = number
	result := card.SetName
	result, _ = strings.CutPrefix(card.SetName, "Pokemon ")
	card.SetName = result

	fmt.Printf("CardName: %s \n CardNumber: %s \n SetNumber: %s", card.Name, card.Number, card.SetNumber)

	psaGrade, err := h.PSA.GetByCardNumberAndSetNumber(card.Number, card.SetNumber)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	bgsGrade, err := h.BGS.GetByCardNumberAndSetNumber(card.Number, card.SetNumber)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	tagGrade, err := h.TAG.GetByCardNumberAndSetNumber(card.Number, card.SetNumber)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	cgcGrade, err := h.CGC.GetByCardNumberAndSetNumber(card.Number, card.SetNumber)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	cardDetail := transformer.ToCardDetail(card, rarity, increasedPrice, percentChange, psaGrade, bgsGrade, tagGrade, cgcGrade)

	return responseJSON(c, echo.Map{
		"detail": cardDetail,
	})
}
