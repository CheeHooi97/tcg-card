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

func (h *Handler) InsertCGC(c echo.Context) error {
	var i struct {
		Url string `json:"url" validate:"required"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	fmt.Println("Start")

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
				cgc.CardName = strings.TrimSpace(card.CardName[:idx+1])
				cgc.Rarity = strings.TrimSpace(card.CardName[idx+1:])
				cgc.CardNumber = card.CardNumber
				parts := strings.Split(card.CardNumber, "/")
				if len(parts) == 2 {
					cgc.CardNumber = parts[0]
					cgc.SetNumber = parts[1]
				}

				check := h.CGC.CheckCardNameAndCardNumberAndSetNameAndRarity(cgc.CardName, cgc.CardNumber, cgc.SetName, cgc.Rarity)

				if !check {
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
				}
			}

			time.Sleep(5 * time.Second)
		}
		time.Sleep(5 * time.Second)
	}

	return responseJSON(c, true)
}

type CGCListData struct {
	URL string `json:"url"`
}

func CGCScrapList(url string) ([]CGCListData, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	var results []CGCListData

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.ccg-cards`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`
			(() => {
				const cards = Array.from(document.querySelectorAll('.card.ng-scope'));
				const baseUrl = "https://www.cgccards.com";
				
				return cards.map(card => {
					const anchor = card.querySelector('a');
					return {
						url: anchor ? baseUrl + anchor.getAttribute('href') : "",
					};
				}).filter(item => item.link !== "");
			})()
		`, &results),
	)

	return results, err
}

type CGCSetData struct {
	SetName string `json:"setName"`
	SetUrl  string `json:"setUrl"`
}

func CGCScrapSets(url string) ([]CGCSetData, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var allResults []CGCSetData

	// 1. Initial Navigation
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`tr.ccg-setcounts-table__row`, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	seenPages := make(map[string]bool)

	for {
		var currentPageData []CGCSetData
		var hasNext bool

		err = chromedp.Run(ctx,
			chromedp.Sleep(1500*time.Millisecond),
			chromedp.Evaluate(`
            (() => {
                const baseUrl = "https://www.cgccards.com";
                const rows = Array.from(document.querySelectorAll('tr.ccg-setcounts-table__row'));

                const results = rows.map(row => {
                    const anchor = row.querySelector('.ccg-setcounts-table__name a');
                    return {
                        setUrl: anchor ? baseUrl + anchor.getAttribute('href') : "",
                        setName: anchor ? anchor.innerText.trim() : ""
                    };
                }).filter(r => r.setUrl !== "");

                const nextBtn = document.querySelector('a.ccg-pager-next');
                const canGoNext = nextBtn && !nextBtn.classList.contains('disabled');

                return {
                    data: results,
                    canGoNext: canGoNext
                };
            })()
        `, &struct {
				Data      *[]CGCSetData `json:"data"`
				CanGoNext *bool         `json:"canGoNext"`
			}{&currentPageData, &hasNext}),
		)
		if err != nil {
			return nil, err
		}

		if len(currentPageData) == 0 {
			fmt.Println("No rows found, stopping.")
			break
		}

		// 🔑 Detect repeated page
		pageKey := currentPageData[0].SetUrl
		if seenPages[pageKey] {
			fmt.Println("Detected repeated page, stopping pagination.")
			break
		}
		seenPages[pageKey] = true

		allResults = append(allResults, currentPageData...)

		if !hasNext {
			fmt.Println("No more pages found.")
			break
		} else {
			fmt.Println("Clicking next page...")
			err = chromedp.Run(ctx,
				chromedp.Click(`a.ccg-pager-next`, chromedp.ByQuery),
				chromedp.Sleep(1200*time.Millisecond),
			)
			if err != nil {
				break
			}
		}
	}

	return allResults, nil
}

type GradeCount struct {
	Grade string `json:"grade"`
	Count string `json:"count"`
}

type CGCCardFullData struct {
	SetName     string       `json:"setName"`
	CardNumber  string       `json:"cardNumber"`
	CardName    string       `json:"cardName"`
	TotalGraded string       `json:"totalGraded"`
	Grades      []GradeCount `json:"grades"`
}

func CGCScrapCards(url string) ([]CGCCardFullData, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var allResults []CGCCardFullData

	// 1. Initial Navigation
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// Wait for the main table rows to be present
		chromedp.WaitVisible(`tr.needs-alignment`, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	maxPages := 15
	pageCount := 0

	for {
		var currentPageResults []CGCCardFullData
		var hasNext bool

		// 2. Scrape Page Content
		err = chromedp.Run(ctx,
			chromedp.Sleep(3*time.Second), // Wait for Angular to populate the tables
			chromedp.Evaluate(`
                (() => {
                    // --- HEADER EXTRACTION (Set Info) ---
                    const headerTitle = document.querySelector('.card-list-header__title');
                    const setName = headerTitle ? headerTitle.innerText.trim() : "";

                    const headerTotalCount = document.querySelector('.card-list-header-count__number');

                    // --- GRADE LABELS ---
                    // Target headers specifically inside the scroller to align with dataColumns
                    const headerCells = Array.from(document.querySelectorAll('#tableScroller thead th.ng-binding'));
                    const gradeLabels = headerCells.map(h => h.innerText.trim()).filter(t => t !== "");

                    // --- TABLE ROWS ---
                    const pinnedRows = Array.from(document.querySelectorAll('.pinned tbody tr.needs-alignment'));
                    const dataRows = Array.from(document.querySelectorAll('#tableScroller tbody tr.needs-alignment.ng-scope'));

                    const data = pinnedRows.map((pRow, index) => {
                        const numCell = pRow.querySelector('.card-list__cardNumber');
                        const nameCell = pRow.querySelector('.card-list__name');
                        const totalCell = pRow.querySelector('.card-list__totalGraded');
                        
                        const dRow = dataRows[index];
                        if (!dRow) return null;

                        // Get Grade Spans (td[0] is the row total, slice(1) gets the grades)
                        const allTds = Array.from(dRow.querySelectorAll('td'));

						// td[0] = total graded
						const gradeStartIndex = 1;

						const grades = gradeLabels.map((label, i) => {
    					const tdIndex = gradeStartIndex + i;
    					const cell = allTds[tdIndex];
    					const span = cell ? cell.querySelector('span.ng-binding') : null;
    					const val = span ? span.innerText.trim() : "";

    					return {
        					grade: label,
        					count: val === "" ? "0" : val
    					};
						});

                        return {
                            setName: setName,
                            cardNumber: numCell ? numCell.innerText.split('\n')[0].trim() : "",
                            cardName: nameCell ? nameCell.innerText.replace(/\s+/g, ' ').trim() : "",
                            totalGraded: totalCell ? totalCell.innerText.trim() : "0",
                            grades: grades
                        };
                    }).filter(r => r !== null && r.cardName !== "");

                    // --- PAGING LOGIC ---
                    const nextBtn = document.querySelector('a.ccg-pager-next');
                    const isNextDisabled = !nextBtn || nextBtn.classList.contains('disabled') || nextBtn.hasAttribute('disabled');

                    return {
                        results: data,
                        canGoNext: !isNextDisabled
                    };
                })()
            `, &struct {
				Results   *[]CGCCardFullData `json:"results"`
				CanGoNext *bool              `json:"canGoNext"`
			}{&currentPageResults, &hasNext}),
		)

		if err != nil {
			return nil, err
		}

		allResults = append(allResults, currentPageResults...)

		// 3. Exit if no more pages
		if !hasNext {
			fmt.Println("No more pages found or single page reached.")
			break
		}

		// 4. Click Next and Repeat
		fmt.Printf("Scraped %d cards. Moving to next page...\n", len(allResults))

		pageCount++
		if pageCount > maxPages {
			fmt.Println("Max page limit reached, stopping pagination.")
			break
		}

		err = chromedp.Run(ctx,
			chromedp.ScrollIntoView(`a.ccg-pager-next`, chromedp.ByQuery),
			chromedp.Click(`a.ccg-pager-next`, chromedp.ByQuery),
			// Wait for the table to refresh with new data
			chromedp.Sleep(1*time.Second),
		)
		if err != nil {
			fmt.Printf("Paging click failed (maybe last page?): %v\n", err)
			break
		}
	}

	fmt.Printf("Total Scraped %d cards....\n", len(allResults))

	return allResults, nil
}
