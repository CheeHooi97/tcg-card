package handler

import (
	"fmt"
	"pkm/errcode"
	"pkm/model"
	"pkm/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CardUpdate(c echo.Context) error {
	var i struct {
		// Id string `json:"id" validate:"required"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	cards, _ := h.BGS.GetAllCards()
	filterCards := make([]*model.BGS, 0)
	for _, card := range cards {
		title := strings.TrimSpace(card.CardName)
		if strings.Contains(title, "FULL ART") {
			filterCards = append(filterCards, card)
		}
	}

	fmt.Println(len(filterCards))

	for _, card := range filterCards {
		title := strings.TrimSpace(card.CardName)

		// if strings.Contains(title, "FULL ART") {
		parts := strings.SplitN(title, "FULL ART", 2)
		card.CardName = strings.TrimSpace(parts[0])
		// }

		card.SetName = card.Description

		fmt.Println("Title: ", title)
		fmt.Println("CardName: ", card.CardName)
		fmt.Println("Rarity: ", card.Rarity)
		fmt.Println("SetName: ", card.SetName)

		if err := h.BGS.Update(card); err != nil {
			return responseError(c, errcode.InternalServerError)
		}
	}

	return responseJSON(c, true)
}
