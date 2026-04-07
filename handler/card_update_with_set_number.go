package handler

import (
	"fmt"
	"pkm/config/set"
	"pkm/errcode"
	"pkm/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CardUpdateSetNumber(c echo.Context) error {
	var i struct {
		// Id string `json:"id" validate:"required"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	cards, err := h.Card.GetAllCards()
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	for _, card := range cards {
		rename, _ := strings.CutPrefix(card.SetName, "Pokemon ")

		// bgs, _ := h.BGS.GetByCardNumberAndSetNumber()
		for checkSet := range set.SetName {
			if rename == checkSet {
				card.SetNumber = set.SetName[checkSet]
				fmt.Println("Card Set Number: ", card.SetNumber)
			}
		}

	}

	return responseJSON(c, cards)
}
