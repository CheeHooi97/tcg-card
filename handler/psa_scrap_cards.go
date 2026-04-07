package handler

import (
	"context"
	"fmt"
	"pkm/errcode"
	"pkm/model"
	"pkm/utils"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) InsertPSA(c echo.Context) error {
	var i struct {
		Urls []string `json:"urls" validate:"required"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	for a, url := range i.Urls {
		fmt.Println("Url: ", a+1)

		sets, err := PSAScrapYear(url)
		if err != nil {
			fmt.Println(err)
			return responseError(c, errcode.InternalServerError)
		}

		time.Sleep(5 * time.Second)

		fmt.Println("Raw set: ", len(sets))

		var newSets []PSASet

		for _, set := range sets {
			if strings.Contains(set.Name, "Pokemon") {
				newSets = append(newSets, set)
			}
		}

		for xxx, newSet := range newSets {
			fmt.Println("Length: ", len(newSets))
			fmt.Println("Current: ", xxx+1)
			fmt.Println("Set: ", newSet.Name)

			cards, _ := PSAScrapSet(newSet.Link)

			if len(cards) != 0 {
				for ccc, card := range cards {
					fmt.Printf("%d of %d cards \n", ccc+1, len(cards))
					check := h.PSA.GetByCardNameAndCardNumberAndDescriptionAndSet(card.CardName, card.CardNumber, card.Description, newSet.Name)

					if !check {
						psa := model.NewPSA()
						psa.CardNumber = card.CardNumber
						psa.CardName = card.CardName
						psa.Description = card.Description
						psa.SetName = newSet.Name
						psa.Total = card.Total
						psa.Grade1 = card.Grade1
						psa.Grade2 = card.Grade2
						psa.Grade3 = card.Grade3
						psa.Grade4 = card.Grade4
						psa.Grade5 = card.Grade5
						psa.Grade6 = card.Grade6
						psa.Grade7 = card.Grade7
						psa.Grade8 = card.Grade8
						psa.Grade9 = card.Grade9
						psa.Grade10 = card.Grade10

						if err := h.PSA.Create(psa); err != nil {
							return responseError(c, errcode.InternalServerError)
						}

						urlCheck := h.PSAUrl.GetByPath(newSet.Link)

						if !urlCheck {
							psaUrl := model.NewPSAUrl()
							psaUrl.SetName = newSet.Name
							psaUrl.Url = newSet.Link

							if err := h.PSAUrl.Create(psaUrl); err != nil {
								return responseError(c, errcode.InternalServerError)
							}
						}
					}
					// else {
					// 	card, err := h.PSA.GetDetailByCardNameAndCardNumberAndDescription(card.CardName, card.CardNumber, card.Description)
					// 	if err != nil {
					// 		return responseError(c, errcode.InternalServerError)
					// 	}

					// 	card.SetName = newSet.Name
					// 	if err := h.PSA.Update(card); err != nil {
					// 		return responseError(c, errcode.InternalServerError)
					// 	}
					// }
				}
			} else {
				continue
			}

			time.Sleep(5 * time.Second)
		}

		time.Sleep(5 * time.Second)
	}

	return responseJSON(c, true)
}

func PSAScrapYear(url string) ([]PSASet, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		// "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"
		// "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var allResults []PSASet

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// WAIT: Let the page load. If you see a Cloudflare checkbox, click it manually!
		chromedp.WaitVisible(`#tableSets`, chromedp.ByID),
		chromedp.Sleep(2*time.Second),

		// FORCE the 50 records view
		chromedp.Evaluate(`
            (() => {
                const select = document.querySelector('select[name="tableSets_length"]');
                if (select) {
                    select.value = "50";
                    // Trigger both 'change' and 'input' events to be sure DataTables reacts
                    select.dispatchEvent(new Event('change', { bubbles: true }));
                    select.dispatchEvent(new Event('input', { bubbles: true }));
                }
            })()
        `, nil),
		chromedp.Sleep(5*time.Second), // Wait for table to reload
	)

	if err != nil {
		return nil, err
	}

	// 2. Start Pagination Loop
	for {
		var pageResults []PSASet

		err := chromedp.Run(ctx,
			chromedp.WaitVisible(`#tableSets`, chromedp.ByID),
			chromedp.Evaluate(`
            (() => {
                const base = "https://www.psacard.com";
                return Array.from(document.querySelectorAll('#tableSets tbody tr')).map(tr => {
                    const anchor = tr.querySelector('td.text-left a:not([href="#"])');
                    return anchor ? { 
                        name: anchor.innerText.trim(), 
                        link: base + anchor.getAttribute('href') 
                    } : null;
                }).filter(i => i !== null);
            })()
        `, &pageResults),
		)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, pageResults...)

		// ✅ Reliable pagination check
		var hasNext bool
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`
            (() => {
                const table = $('#tableSets').DataTable();
                const info = table.page.info();
                return info.page + 1 < info.pages;
            })()
        `, &hasNext),
		)
		if err != nil {
			return nil, err
		}

		if !hasNext {
			fmt.Println("Final page reached.")
			break
		}

		fmt.Println("Moving to next page...")
		err = chromedp.Run(ctx,
			chromedp.Click(`#tableSets_next`, chromedp.ByID),
			chromedp.Sleep(1500*time.Millisecond),
		)
		if err != nil {
			return nil, err
		}
	}

	return allResults, nil
}

func PSAScrapSet(url string) ([]PSACardData, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		// Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0
		// Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var results []PSACardData

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#tablePSA`, chromedp.ByID),
		chromedp.Sleep(2*time.Second),

		// 1. Force 500 records
		chromedp.Evaluate(`
        (() => {
            const select = document.querySelector('select[name="tablePSA_length"]');
            if (select) {
                select.value = "500";
                select.dispatchEvent(new Event('change', { bubbles: true }));
            }
        })()
    `, nil),
		chromedp.Sleep(5*time.Second),

		// 2. Scrape with Grades
		chromedp.Evaluate(`
        (() => {
            const rows = Array.from(document.querySelectorAll('#tablePSA tbody tr[role="row"]'));
            
            return rows.map(tr => {
                const cells = Array.from(tr.querySelectorAll('td'));
                
                if (cells.length < 16 || cells[2].innerText.includes("TOTAL POPULATION")) {
                    return null;
                }

                // Helper to get the top number from the Grade div stack
                const getGradeVal = (cell) => {
                    const firstDiv = cell.querySelector('div');
                    return firstDiv ? firstDiv.innerText.trim() : "0";
                };

                const nameElem = cells[2].querySelector('strong');
                const cardName = nameElem ? nameElem.innerText.trim() : "";
                
                let cellClone = cells[2].cloneNode(true);
                const link = cellClone.querySelector('a');
                const strong = cellClone.querySelector('strong');
                if (link) link.remove();
                if (strong) strong.remove();
                const description = cellClone.innerText.trim();

                return {
                    cardNumber:  cells[1].innerText.trim(),
                    cardName:    cardName,
                    description: description,
                    grade6:      getGradeVal(cells[11]),
                    grade7:      getGradeVal(cells[12]),
                    grade8:      getGradeVal(cells[13]),
                    grade9:      getGradeVal(cells[14]),
                    grade10:     getGradeVal(cells[15]),
                    total:       getGradeVal(cells[16])
                };
            }).filter(i => i !== null);
        })()
    `, &results),
	)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Scraped %d cards from the set.\n", len(results))
	return results, nil
}

type PSASet struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type PSACardData struct {
	CardNumber  string `json:"cardNumber"`
	CardName    string `json:"cardName"`
	Description string `json:"description"`
	Grade1      string `json:"Grade1"`
	Grade2      string `json:"Grade2"`
	Grade3      string `json:"Grade3"`
	Grade4      string `json:"Grade4"`
	Grade5      string `json:"Grade5"`
	Grade6      string `json:"Grade6"`
	Grade7      string `json:"Grade7"`
	Grade8      string `json:"Grade8"`
	Grade9      string `json:"Grade9"`
	Grade10     string `json:"Grade10"`
	Total       string `json:"total"`
}
