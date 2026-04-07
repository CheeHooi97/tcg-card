package handler

import (
	"fmt"
	"pkm/errcode"
	"pkm/model"
	"pkm/utils"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CheckCGC(c echo.Context) error {
	var i struct {
		Url string `json:"url" validate:"required"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	lists, err := CGCScrapList(i.Url)
	if err != nil {
		fmt.Println(err)
		return responseError(c, errcode.InternalServerError)
	}

	time.Sleep(3 * time.Second)

	for z, list := range lists {
		fmt.Printf("%d/%d List \n", z+1, len(lists))

		sets, err := CGCScrapSets(list.URL)
		if err != nil {
			return responseError(c, errcode.InternalServerError)
		}

		time.Sleep(3 * time.Second)

		for y, set := range sets {
			fmt.Printf("%d/%d Set \n", y+1, len(sets))
			fmt.Println("Set: ", set.SetName)

			cards, err := CGCScrapCards(set.SetUrl)
			if err != nil {
				return responseError(c, errcode.InternalServerError)
			}

			for x, card := range cards {
				fmt.Printf("Set %d: %d of %d cards \n", y+1, x+1, len(cards))

				cgc := model.NewCGC()
				cgc.SetName = card.SetName
				cgc.Total = card.TotalGraded

				idx := strings.Index(card.CardName, ")")
				if idx == -1 {
					cgc.CardName = strings.TrimSpace(card.CardName)
					cgc.Rarity = ""
				} else {
					cgc.CardName = strings.TrimSpace(card.CardName[:idx+1])
					cgc.Rarity = strings.TrimSpace(card.CardName[idx+1:])
				}

				cgc.CardNumber = card.CardNumber
				parts := strings.Split(card.CardNumber, "/")
				if len(parts) == 2 {
					cgc.CardNumber = strings.TrimSpace(parts[0])
					cgc.SetNumber = strings.TrimSpace(parts[1])
				}

				exists := h.CGC.CheckCardNameAndCardNumberAndSetNameAndRarity(cgc.CardName, cgc.CardNumber, cgc.SetName, cgc.Rarity)

				if !exists {
					for _, grade := range card.Grades {
						switch grade.Grade {
						case "1":
							cgc.Grade1 = grade.Count
						case "1.5":
							cgc.Grade1_5 = grade.Count
						case "2":
							cgc.Grade2 = grade.Count
						case "2.5":
							cgc.Grade2_5 = grade.Count
						case "3":
							cgc.Grade3 = grade.Count
						case "3.5":
							cgc.Grade3_5 = grade.Count
						case "4":
							cgc.Grade4 = grade.Count
						case "4.5":
							cgc.Grade4_5 = grade.Count
						case "5":
							cgc.Grade5 = grade.Count
						case "5.5":
							cgc.Grade5_5 = grade.Count
						case "6":
							cgc.Grade6 = grade.Count
						case "6.5":
							cgc.Grade6_5 = grade.Count
						case "7":
							cgc.Grade7 = grade.Count
						case "7.5":
							cgc.Grade7_5 = grade.Count
						case "8":
							cgc.Grade8 = grade.Count
						case "8.5":
							cgc.Grade8_5 = grade.Count
						case "9":
							cgc.Grade9 = grade.Count
						case "Mint+ 9.5":
							cgc.Grade9_5 = grade.Count
						case "Gem Mint 10":
							cgc.Grade10 = grade.Count
						case "Pristine 10":
							cgc.Grade10P = grade.Count
						case "AU":
							continue
						case "AA":
							continue
						case "Perfect 10":
							continue
						default:
							continue
						}
					}

					if err := h.CGC.Create(cgc); err != nil {
						return responseError(c, errcode.InternalServerError)
					}

					continue
				}

				dbCGC, err := h.CGC.GetByCardNameAndCardNumberAndSetNameAndRarity(cgc.CardName, cgc.CardNumber, cgc.SetName, cgc.Rarity)
				if err != nil {
					return responseError(c, errcode.InternalServerError)
				}

				logging := model.NewCGCLogging()
				logging.CGCId = dbCGC.Id
				logging.CardName = cgc.CardName
				logging.CardNumber = cgc.CardNumber
				logging.SetNumber = cgc.SetNumber
				logging.SetName = cgc.SetName
				logging.Rarity = cgc.Rarity
				logging.Total = cgc.Total

				for _, grade := range card.Grades {
					switch grade.Grade {
					case "1":
						logging.Grade1 = grade.Count
					case "1.5":
						logging.Grade1_5 = grade.Count
					case "2":
						logging.Grade2 = grade.Count
					case "2.5":
						logging.Grade2_5 = grade.Count
					case "3":
						logging.Grade3 = grade.Count
					case "3.5":
						logging.Grade3_5 = grade.Count
					case "4":
						logging.Grade4 = grade.Count
					case "4.5":
						logging.Grade4_5 = grade.Count
					case "5":
						logging.Grade5 = grade.Count
					case "5.5":
						logging.Grade5_5 = grade.Count
					case "6":
						logging.Grade6 = grade.Count
					case "6.5":
						logging.Grade6_5 = grade.Count
					case "7":
						logging.Grade7 = grade.Count
					case "7.5":
						logging.Grade7_5 = grade.Count
					case "8":
						logging.Grade8 = grade.Count
					case "8.5":
						logging.Grade8_5 = grade.Count
					case "9":
						logging.Grade9 = grade.Count
					case "Mint+ 9.5":
						logging.Grade9_5 = grade.Count
					case "Gem Mint 10":
						logging.Grade10 = grade.Count
					case "Pristine 10":
						logging.Grade10P = grade.Count
					case "AU":
						continue
					case "AA":
						continue
					case "Perfect 10":
						continue
					default:
						continue
					}
				}

				if err := h.CGCLogging.Create(logging); err != nil {
					return responseError(c, errcode.InternalServerError)
				}
			}

			time.Sleep(5 * time.Second)
		}
		time.Sleep(5 * time.Second)
	}

	return responseJSON(c, true)
}
