package handler

import (
	"fmt"
	"pkm/config/set"
	"pkm/model"
	"pkm/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateEachCardBasedOnSetName(c echo.Context) error {
	var i struct {
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	cardLists, _ := h.Card.GetAllCards()

	fmt.Println("len of cards list: ", len(cardLists))

	var length int
	newCards := make([]*model.Card, 0)
	for _, card := range cardLists {
		if card.SetNumber != "" {
			newCards = append(newCards, card)
			continue
		}

		for setName, setNumber := range set.SetName {
			if card.SetName == "Pokemon "+setName {
				if card.SetNumber == "" {
					card.SetNumber = setNumber

					if strings.Contains(card.SetName, "Pokemon Neo") || strings.Contains(card.SetName, "Pokemon Gym Challenge") ||
						strings.Contains(card.SetName, "Pokemon Gym Heroes") || strings.Contains(card.SetName, "Pokemon Team Rocket") ||
						strings.Contains(card.SetName, "Pokemon Fossil") || strings.Contains(card.SetName, "Pokemon Jungle") ||
						strings.Contains(card.SetName, "Pokemon Base Set") {
						if strings.Contains(card.Name, "1st Edition") {
							card.SetName = card.SetName + " 1st Edition"
						} else if strings.Contains(card.Name, "Shadowless") {
							card.SetName = card.SetName + " Shadowless"
						} else {
							card.SetName = card.SetName + " Unlimited"
						}
					}

					length++
					// if err := h.Card.Update(card); err != nil {
					// 	return responseError(c, errcode.InternalServerError)
					// }

					newCards = append(newCards, card)
					break
				}
			}
		}
	}

	fmt.Println("len:", length)

	return responseJSON(c, echo.Map{
		"lists": newCards,
	})
}
