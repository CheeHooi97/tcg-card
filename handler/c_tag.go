package handler

import (
	"context"
	"fmt"
	"pkm/config/rarity"
	"pkm/model"
	"pkm/utils"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) InsertTAG(c echo.Context) error {
	var i struct {
		Urls []string `json:"urls" validate:"required"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	for _, yearPath := range i.Urls {
		var setUrlLists []TAGYearSet
		err := retryWithBackoff(fmt.Sprintf("tag scrape year %s", yearPath), 3, 3*time.Second, func() error {
			var scrapeErr error
			setUrlLists, scrapeErr = TAGScrapSet(yearPath)
			return scrapeErr
		})
		if err != nil {
			fmt.Println("TAGScrapSet error:", err)
			continue
		}

		time.Sleep(3 * time.Second)

		for y, setUrl := range setUrlLists {
			// if y < 42 {
			// 	continue
			// }

			var cardLists TAGPopResponse
			err := retryWithBackoff(fmt.Sprintf("tag scrape cards set %s", setUrl.Name), 3, 3*time.Second, func() error {
				var scrapeErr error
				cardLists, scrapeErr = TAGScrapCards(setUrl.Link)
				if scrapeErr != nil {
					return scrapeErr
				}
				if len(cardLists.Cards) == 0 {
					return fmt.Errorf("empty cards")
				}
				return nil
			})
			if err != nil {
				fmt.Printf("Skip TAG set %s after retries: %v\n", setUrl.Name, err)
				continue
			}

			for x, card := range cardLists.Cards {
				if x == 0 {
					fmt.Printf("%d of %d Lists \n", y+1, len(setUrlLists))
					fmt.Println("Set Name: ", setUrl.Name)
				}

				fmt.Printf("%d of %d Cards \n", x+1, len(cardLists.Cards))

				name, rarity := utils.SplitNameAndRarity(card.CardName, rarity.TAGRarities)
				parts := strings.Split(card.CardNumber, "/")
				if len(parts) == 2 {
					card.CardNumber = strings.TrimSpace(parts[0])
				} else {
					// fallback: keep original or handle special cases
					card.CardNumber = strings.TrimSpace(card.CardNumber)
				}

				result := cardLists.SetName
				if v, ok := strings.CutPrefix(cardLists.SetName, "Pokémon "); ok {
					result = v
				}

				if strings.Contains(name, "FULL ART -") {
					card.CardName = strings.TrimSpace(
						strings.TrimSuffix(name, "FULL ART -"),
					)
				} else if strings.Contains(name, "Full Art -") {
					card.CardName = strings.TrimSpace(
						strings.TrimSuffix(name, "Full Art -"),
					)
				}

				check := h.TAG.GetByCardNameAndCardNumberAndSetName(card.CardName, card.CardNumber, cardLists.SetName)

				if !check {
					tag := model.NewTAG()
					tag.CardSet = setUrl.Name
					tag.Total = card.Total
					tag.GradeVA = card.GradeVA
					tag.Grade1 = card.Grade1
					tag.Grade1_5 = card.Grade1_5
					tag.Grade2 = card.Grade2
					tag.Grade2_5 = card.Grade2_5
					tag.Grade3 = card.Grade3
					tag.Grade3_5 = card.Grade3_5
					tag.Grade4 = card.Grade4
					tag.Grade4_5 = card.Grade4_5
					tag.Grade5 = card.Grade5
					tag.Grade5_5 = card.Grade5_5
					tag.Grade6 = card.Grade6
					tag.Grade6_5 = card.Grade6_5
					tag.Grade7 = card.Grade7
					tag.Grade7_5 = card.Grade7_5
					tag.Grade8 = card.Grade8
					tag.Grade8_5 = card.Grade8_5
					tag.Grade9 = card.Grade9
					tag.Grade10 = card.Grade10
					tag.Grade10P = card.Grade10P

					tag.CardName = name
					tag.Rarity = rarity

					tag.SetName = result
					if len(parts) == 2 {
						tag.CardNumber = strings.TrimSpace(parts[0])
						tag.SetNumber = strings.TrimSpace(parts[1])
					} else {
						// fallback: keep original or handle special cases
						tag.CardNumber = strings.TrimSpace(card.CardNumber)
						tag.SetNumber = ""
					}

					if err := h.TAG.Create(tag); err != nil {
						fmt.Printf("TAG create error, skip card %s %s: %v\n", tag.CardNumber, tag.CardName, err)
						continue
					}

					logging := model.NewTAGLogging()
					logging.TagId = tag.Id
					logging.CardName = name
					logging.CardNumber = card.CardNumber
					logging.SetName = result
					logging.CardSet = setUrl.Name
					logging.Rarity = rarity
					logging.Total = card.Total
					logging.GradeVA = card.GradeVA
					logging.Grade1 = card.Grade1
					logging.Grade1_5 = card.Grade1_5
					logging.Grade2 = card.Grade2
					logging.Grade2_5 = card.Grade2_5
					logging.Grade3 = card.Grade3
					logging.Grade3_5 = card.Grade3_5
					logging.Grade4 = card.Grade4
					logging.Grade4_5 = card.Grade4_5
					logging.Grade5 = card.Grade5
					logging.Grade5_5 = card.Grade5_5
					logging.Grade6 = card.Grade6
					logging.Grade6_5 = card.Grade6_5
					logging.Grade7 = card.Grade7
					logging.Grade7_5 = card.Grade7_5
					logging.Grade8 = card.Grade8
					logging.Grade8_5 = card.Grade8_5
					logging.Grade9 = card.Grade9
					logging.Grade10 = card.Grade10
					logging.Grade10P = card.Grade10P

					if len(parts) == 2 {
						logging.SetNumber = strings.TrimSpace(parts[1])
					}

					if err := h.TAGLogging.Create(logging); err != nil {
						fmt.Printf("TAG logging create error, skip card %s %s: %v\n", logging.CardNumber, logging.CardName, err)
						continue
					}

					urlCheck := h.TAGUrl.GetByPath(setUrl.Link)

					if !urlCheck {
						tagUrl := model.NewTAGUrl()
						tagUrl.Url = setUrl.Link

						if err := h.TAGUrl.Create(tagUrl); err != nil {
							fmt.Printf("TAG url create error, continue: %v\n", err)
							continue
						}
					}
				} else {
					dbTag, err := h.TAG.GetDetailByCardNameAndCardNumberAndSetName(name, card.CardNumber, setUrl.Name)
					if err != nil {
						fmt.Printf("TAG detail lookup error, skip card %s %s: %v\n", card.CardNumber, name, err)
						continue
					}

					logging := model.NewTAGLogging()
					if dbTag != nil && dbTag.Id != "" {
						logging.TagId = dbTag.Id
					}

					logging.CardName = name
					logging.CardNumber = card.CardNumber
					logging.SetName = result
					logging.CardSet = setUrl.Name
					logging.Rarity = rarity
					logging.Total = card.Total
					logging.GradeVA = card.GradeVA
					logging.Grade1 = card.Grade1
					logging.Grade1_5 = card.Grade1_5
					logging.Grade2 = card.Grade2
					logging.Grade2_5 = card.Grade2_5
					logging.Grade3 = card.Grade3
					logging.Grade3_5 = card.Grade3_5
					logging.Grade4 = card.Grade4
					logging.Grade4_5 = card.Grade4_5
					logging.Grade5 = card.Grade5
					logging.Grade5_5 = card.Grade5_5
					logging.Grade6 = card.Grade6
					logging.Grade6_5 = card.Grade6_5
					logging.Grade7 = card.Grade7
					logging.Grade7_5 = card.Grade7_5
					logging.Grade8 = card.Grade8
					logging.Grade8_5 = card.Grade8_5
					logging.Grade9 = card.Grade9
					logging.Grade10 = card.Grade10
					logging.Grade10P = card.Grade10P

					if len(parts) == 2 {
						logging.SetNumber = strings.TrimSpace(parts[1])
					}

					if dbTag != nil && dbTag.Id != "" {
						if dbTag.SetName != "" {
							logging.SetName = dbTag.SetName
						}
						if dbTag.CardSet != "" {
							logging.CardSet = dbTag.CardSet
						}
						if dbTag.Rarity != "" {
							logging.Rarity = dbTag.Rarity
						}
						if dbTag.SetNumber != "" {
							logging.SetNumber = dbTag.SetNumber
						}
					}

					if err := h.TAGLogging.Create(logging); err != nil {
						fmt.Printf("TAG logging create error, skip card %s %s: %v\n", logging.CardNumber, logging.CardName, err)
						continue
					}
				}

				if len(i.Urls) == 0 {
					break
				}
			}

			time.Sleep(3 * time.Second)
		}
	}

	fmt.Println("------------------------------")
	fmt.Println("End of scrapping")

	return responseJSON(c, true)
}

func TAGScrapCards(url string) (TAGPopResponse, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var results TAGPopResponse

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// Wait for the table rows to actually exist
		chromedp.WaitVisible(`tr.MuiTableRow-root`, chromedp.ByQuery),
		// Crucial: The page needs extra time to fill the numbers into the links
		chromedp.Sleep(5*time.Second),
		chromedp.Evaluate(`
			(() => {
			 	const urlParams = new URLSearchParams(window.location.search);
                const setName = urlParams.get('setName') || document.title;

				const rows = Array.from(document.querySelectorAll('tbody tr.MuiTableRow-root'));
				
				const cards = rows.map(row => {
					const cells = Array.from(row.querySelectorAll('td'));
					if (cells.length < 15) return null;

					return {
						cardNumber: cells[0].innerText.trim(),
						cardName:   cells[1].innerText.replace(/\n/g, ' ').trim(),
						gradeVA :   cells[cells.length - 21].innerText.trim() || "0",
						grade1  :   cells[cells.length - 20].innerText.trim() || "0",
						grade1_5:   cells[cells.length - 19].innerText.trim() || "0",
						grade2  :   cells[cells.length - 18].innerText.trim() || "0",
						grade2_5:   cells[cells.length - 17].innerText.trim() || "0",
						grade3  :   cells[cells.length - 16].innerText.trim() || "0",
						grade3_5:   cells[cells.length - 15].innerText.trim() || "0",
						grade4:     cells[cells.length - 14].innerText.trim() || "0",
						grade4_5:   cells[cells.length - 13].innerText.trim() || "0",
						grade5:     cells[cells.length - 12].innerText.trim() || "0",
						grade5_5:   cells[cells.length - 11].innerText.trim() || "0",
						grade6:     cells[cells.length - 10].innerText.trim() || "0",
						grade6_5:   cells[cells.length - 9].innerText.trim() || "0",
						grade7:     cells[cells.length - 8].innerText.trim() || "0",
						grade7_5:   cells[cells.length - 7].innerText.trim() || "0",
						grade8:     cells[cells.length - 6].innerText.trim() || "0",
						grade8_5:   cells[cells.length - 5].innerText.trim() || "0",
						grade9:     cells[cells.length - 4].innerText.trim() || "0",
						grade10:    cells[cells.length - 3].innerText.trim() || "0",
						grade10P:   cells[cells.length - 2].innerText.trim() || "0",
						total:      cells[cells.length - 1].innerText.trim() || "0"
					};
				}).filter(item => item !== null && item.cardNumber !== "");
				return {
                    setName: setName,
                    cards: cards
                };
			})()
		`, &results),
	)

	if err != nil {
		return TAGPopResponse{}, err
	}

	fmt.Printf("Scraped %d cards successfully:\n", len(results.Cards))
	for _, res := range results.Cards {
		fmt.Println("res: ", res)

	}

	return results, nil
}

type TAGPopData struct {
	CardName    string `json:"cardName"`
	CardNumber  string `json:"cardNumber"`
	Description string `json:"description"`
	Total       string `json:"total"`
	GradeVA     string `json:"GradeVA"`
	Grade1      string `json:"Grade1"`
	Grade1_5    string `json:"Grade1_5"`
	Grade2      string `json:"Grade2"`
	Grade2_5    string `json:"Grade2_5"`
	Grade3      string `json:"Grade3"`
	Grade3_5    string `json:"Grade3_5"`
	Grade4      string `json:"Grade4"`
	Grade4_5    string `json:"Grade4_5"`
	Grade5      string `json:"Grade5"`
	Grade5_5    string `json:"Grade5_5"`
	Grade6      string `json:"Grade6"`
	Grade6_5    string `json:"Grade6_5"`
	Grade7      string `json:"Grade7"`
	Grade7_5    string `json:"Grade7_5"`
	Grade8      string `json:"Grade8"`
	Grade8_5    string `json:"grade8_5"`
	Grade9      string `json:"grade9"`
	Grade10     string `json:"grade10"`
	Grade10P    string `json:"grade10P"`
}

type TAGPopResponse struct {
	SetName string       `json:"setName"`
	Cards   []TAGPopData `json:"cards"`
}

type TAGYearSet struct {
	Link string `json:"link"`
	Name string `json:"name"`
}

func TAGScrapSet(yearUrl string) ([]TAGYearSet, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/google-chrome"),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("window-size", "1920,1080"),
		chromedp.Flag("headless", true),
		// chromedp.Flag("disable-blink-features", "AutomationControlled"),
		// chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var results []TAGYearSet

	err := chromedp.Run(ctx,
		chromedp.Navigate(yearUrl),
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
	)

	if err != nil {
		return nil, err
	}

	scrapeSetJs := `
		(() => {
			const baseUrl = window.location.origin || "https://my.taggrading.com";
			const seen = new Set();
			const output = [];
			const rowSelectors = [
				'tbody tr.MuiTableRow-root',
				'tr.MuiTableRow-root',
				'table tbody tr',
				'tbody tr'
			];

			const collect = (anchor, fallbackName) => {
				if (!anchor) return;
				const hrefRaw = anchor.getAttribute('href') || '';
				if (!hrefRaw) return;

				const link = new URL(hrefRaw, baseUrl).href;
				if (!link.includes('/pop-report/')) return;
				if (link === window.location.href) return;
				if (seen.has(link)) return;

				const name = (fallbackName || anchor.textContent || '')
					.replace(/\s+/g, ' ')
					.trim();

				seen.add(link);
				output.push({ link, name });
			};

			for (const selector of rowSelectors) {
				const rows = Array.from(document.querySelectorAll(selector));
				rows.forEach(row => {
					const anchor = row.querySelector('a[href*="/pop-report/"]');
					const bold = row.querySelector('b');
					const firstCell = row.querySelector('td');
					const fallbackName = (bold && bold.textContent) || (firstCell && firstCell.textContent) || '';
					collect(anchor, fallbackName);
				});
			}

			// Fallback for layouts where links are not inside table rows.
			if (output.length === 0) {
				const links = Array.from(document.querySelectorAll('a[href*="/pop-report/"]'));
				links.forEach(anchor => collect(anchor, anchor.textContent || ''));
			}

			return output;
		})()
	`

	// Retry several times for SPA hydration / delayed data load.
	for attempt := 0; attempt < 8; attempt++ {
		results = nil
		err = chromedp.Run(ctx, chromedp.Evaluate(scrapeSetJs, &results))
		if err != nil {
			return nil, err
		}

		if len(results) > 0 {
			break
		}

		time.Sleep(1500 * time.Millisecond)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no sets found from %s", yearUrl)
	}

	fmt.Printf("Scraped %d sets successfully:\n", len(results))
	return results, nil
}
