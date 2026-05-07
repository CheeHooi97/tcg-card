package handler

import (
	"pkm/errcode"
	"pkm/middleware"
	"pkm/transformer"
	"pkm/utils"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

var scanCardNumberRegex = regexp.MustCompile(`\b\d{1,3}\s*/\s*\d{1,3}\b`)

func (h *Handler) ScanCard(c echo.Context) error {
	var i struct {
		ScanText string `json:"scanText"`
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

	parsedKeyword := buildScanKeyword(i.ScanText)
	if parsedKeyword == "" {
		return responseJSON(c, echo.Map{
			"scanText":      i.ScanText,
			"parsedKeyword": "",
			"lists":         []any{},
		})
	}

	cards, err := h.Card.SearchCardKeywords(parsedKeyword)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	rarityMap := make(map[string]string, 0)
	increasedPriceMap := make(map[string]string)
	percentChangeMap := make(map[string]string)

	cardLists := transformer.ToCardLists(cards, rarityMap, increasedPriceMap, percentChangeMap)

	return responseJSON(c, echo.Map{
		"scanText":      i.ScanText,
		"parsedKeyword": parsedKeyword,
		"lists":         cardLists,
	})
}

func buildScanKeyword(scanText string) string {
	text := strings.TrimSpace(scanText)
	if text == "" {
		return ""
	}

	clean := strings.Join(strings.Fields(text), " ")
	lower := strings.ToLower(clean)
	tokens := strings.Fields(clean)

	// Keep likely set number pattern such as "58/102".
	cardNo := scanCardNumberRegex.FindString(clean)

	filtered := make([]string, 0, len(tokens))
	for _, t := range tokens {
		token := strings.TrimSpace(t)
		if token == "" {
			continue
		}

		// Skip noisy OCR tokens that hurt matching.
		l := strings.ToLower(strings.Trim(token, ".,:;()[]{}"))
		if l == "pokemon" || l == "card" || l == "trading" || l == "game" || l == "tcg" {
			continue
		}
		filtered = append(filtered, strings.Trim(token, ".,:;()[]{}"))
	}

	base := strings.TrimSpace(strings.Join(filtered, " "))
	if base == "" && cardNo == "" {
		return ""
	}

	if strings.Contains(lower, "pokemon") && !strings.Contains(strings.ToLower(base), "pokemon") {
		base = strings.TrimSpace(base + " pokemon")
	}

	if cardNo != "" && !strings.Contains(base, cardNo) {
		base = strings.TrimSpace(base + " " + cardNo)
	}

	return strings.TrimSpace(base)
}
