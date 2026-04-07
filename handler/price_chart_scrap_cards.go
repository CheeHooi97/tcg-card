package handler

import (
	"context"
	"fmt"
	"pkm/config"
	"pkm/errcode"
	"pkm/kit/oss"
	"pkm/model"
	"pkm/utils"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) PriceChartCards(c echo.Context) error {
	var i struct {
		Url string `json:"url" validate:"required"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	// detail, err := h.Card.GetById("1768799412640396734")
	// if err != nil {
	// 	return responseError(c, errcode.InternalServerError)
	// }

	// url, err := oss.GetSignURL(config.OSSBucket, detail.PhotoUrl)

	// fmt.Println("url:", url)

	// get all sets link
	sets, err := PriceChartScrapSet(i.Url)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	time.Sleep(3 * time.Second)

	// loop each set
	for y, set := range sets {
		fmt.Printf("%d/%d sets \n", y+1, len(sets))

		// if y < 18 {
		// 	continue
		// }

		// if y != 21 {
		// 	continue
		// }

		// get cards list of each set
		cards, err := ScrapPriceChartCards(set.Link)
		if err != nil {
			return responseError(c, errcode.InternalServerError)
		}

		time.Sleep(3 * time.Second)

		// get card detail of each card
		for z, card := range cards {
			fmt.Printf("%d/%d cards \n", z+1, len(cards))

			fmt.Println("card:", card)

			// if z < 200 {
			// 	continue
			// }
			// 420 cards

			// get card detail of each card
			cardDetail, err := ScrapCardDetails(card.Link)
			if err != nil {
				return responseError(c, errcode.InternalServerError)
			}

			fmt.Println("cardDetail:", cardDetail)

			check := h.Card.GetByCardNameAndSet(card.Name, cardDetail.SetName)

			if !check {
				ca := model.NewCard()

				fileByte, fileName, err := oss.ProcessImageUrl(cardDetail.ImageURL, ca.Id)
				if err != nil {
					return responseError(c, errcode.InternalServerError)
				}

				fmt.Println("Name: ", fileName)

				if err := oss.Upload(config.OSSBucket, fileName, fileByte); err != nil {
					return responseError(c, errcode.InternalServerError)
				}

				ca.Name = card.Name
				ca.SetName = cardDetail.SetName
				ca.PhotoUrl = fileName
				ca.Ungrade = card.Price
				ca.Grade7 = cardDetail.Grade7
				ca.Grade8 = cardDetail.Grade8
				ca.Grade9 = cardDetail.Grade9
				ca.Grade9_5 = cardDetail.Grade9_5
				ca.Grade10 = cardDetail.Grade10
				if err := h.Card.Create(ca); err != nil {
					return responseError(c, errcode.InternalServerError)
				}

				cardPrice := model.NewCardPrice()
				cardPrice.CardId = ca.Id
				cardPrice.Name = card.Name
				cardPrice.Set = cardDetail.SetName
				cardPrice.Ungrade = cardDetail.Ungraded
				cardPrice.Grade7 = cardDetail.Grade7
				cardPrice.Grade8 = cardDetail.Grade8
				cardPrice.Grade9 = cardDetail.Grade9
				cardPrice.Grade9_5 = cardDetail.Grade9_5
				cardPrice.Grade10 = cardDetail.Grade10

				if err := h.CardPrice.Create(cardPrice); err != nil {
					return responseError(c, errcode.InternalServerError)
				}
			}
			time.Sleep(3 * time.Second)
		}
		time.Sleep(3 * time.Second)
	}

	return responseJSON(c, true)
}

type PriceChartSet struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

func PriceChartScrapSet(yearUrl string) ([]PriceChartSet, error) {
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

	var results []PriceChartSet

	err := chromedp.Run(ctx,
		// Navigate to the PriceCharting page
		chromedp.Navigate(yearUrl),

		// Wait for the container that holds the list of sets
		chromedp.WaitVisible(`.home-box.all ul`, chromedp.ByQuery),

		// Run JavaScript to extract names and links
		chromedp.Evaluate(`
			(() => {
				// Target the specific container from your HTML snippet
				const container = document.querySelector('.home-box.all');
				if (!container) return [];

				const listItems = Array.from(container.querySelectorAll('ul li'));
				const baseUrl = "https://www.pricecharting.com";

				return listItems.map(li => {
					const anchor = li.querySelector('a');
					if (!anchor) return null;

					return {
						name: anchor.innerText.trim(),
						link: anchor.getAttribute('href').startsWith('http') 
                              ? anchor.getAttribute('href') 
                              : baseUrl + anchor.getAttribute('href')
					};
				}).filter(item => item !== null && item.link !== "");
			})()
		`, &results),
	)

	if err != nil {
		return nil, err
	}

	return results, nil
}

type CardData struct {
	ProductID string `json:"productId"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	Link      string `json:"link"`
}

func ScrapPriceChartCards(setUrl string) ([]CardData, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()

	var results []CardData

	err := chromedp.Run(ctx,
		chromedp.Navigate(setUrl),
		chromedp.WaitVisible(`#games_table tbody tr[data-product]`, chromedp.ByQuery),

		chromedp.ActionFunc(func(ctx context.Context) error {
			seen := make(map[string]bool)
			noNewDataCount := 0
			iteration := 0

			// 1. Detect Total
			var totalExpected int
			chromedp.Evaluate(`(() => {
				const selectors = ['.section-subtitle b', '#console-header b', '.f-12 b'];
				for (let s of selectors) {
					const el = document.querySelector(s);
					if (el && el.innerText.includes('/')) {
						const m = el.innerText.match(/\/ (\d+)/);
						if (m) return parseInt(m[1]);
					}
				}
				return 0;
			})()`, &totalExpected).Do(ctx)

			fmt.Printf("--- Scraping Started ---\nTarget: %d cards\n", totalExpected)

			for {
				iteration++
				var currentBatch []CardData

				// 2. Extract Data
				err := chromedp.Evaluate(`
					(() => {
						const rows = Array.from(document.querySelectorAll('#games_table tbody tr[data-product]'));
						const baseUrl = "https://www.pricecharting.com";
						return rows.map(row => {
							const titleA = row.querySelector('.title a');
							return {
								productId: row.getAttribute('data-product') || "",
								name: titleA ? titleA.innerText.trim() : "",
								price: row.querySelector('.used_price .js-price')?.innerText.trim() || "N/A",
								image: row.querySelector('.photo')?.getAttribute('src') || "",
								link: titleA ? (titleA.getAttribute('href').startsWith('http') ? titleA.getAttribute('href') : baseUrl + titleA.getAttribute('href')) : ""
							};
						});
					})()
				`, &currentBatch).Do(ctx)

				if err != nil {
					return err
				}

				newFound := 0
				for _, card := range currentBatch {
					if card.ProductID != "" && !seen[card.ProductID] {
						seen[card.ProductID] = true
						results = append(results, card)
						newFound++
					}
				}

				fmt.Printf("Iteration %d | Total: %d | New: %d\n", iteration, len(results), newFound)

				if totalExpected > 0 && len(results) >= totalExpected {
					fmt.Println("Target reached successfully.")
					break
				}

				// 3. SLOW SCROLL TRIGGER
				// Instead of one big jump, we do 4 small jumps to ensure we hit the trigger
				for i := 0; i < 4; i++ {
					chromedp.Evaluate(`window.scrollBy(0, 600);`, nil).Do(ctx)
					time.Sleep(500 * time.Millisecond) // Mini-pause between steps
				}

				// Final check: Scroll the very last row into view
				chromedp.Evaluate(`
					const rows = document.querySelectorAll('#games_table tbody tr');
					if (rows.length > 0) {
						rows[rows.length - 1].scrollIntoView({behavior: "smooth", block: "end"});
					}
				`, nil).Do(ctx)

				// 4. STALL LOGIC (Patience increased to 5 for slow loading)
				if newFound == 0 && len(results) > 0 {
					noNewDataCount++
					fmt.Printf("Waiting for new data... (Attempt %d/5)\n", noNewDataCount)
					if noNewDataCount > 4 {
						break
					}
				} else {
					noNewDataCount = 0
				}

				// Long pause to let the site "breathe" and load next items
				time.Sleep(4 * time.Second)
			}
			return nil
		}),
	)

	return results, err
}

type FullCardDetails struct {
	SetName  string `json:"setName"`
	ImageURL string `json:"imageUrl"`
	Ungraded string `json:"ungraded"`
	Grade7   string `json:"grade7"`
	Grade8   string `json:"grade8"`
	Grade9   string `json:"grade9"`
	Grade9_5 string `json:"grade9_5"`
	Grade10  string `json:"grade10"`
}

func ScrapCardDetails(urlLink string) (*FullCardDetails, error) {
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

	var result FullCardDetails

	err := chromedp.Run(ctx,
		chromedp.Navigate(urlLink),
		chromedp.WaitVisible(`#product_details .cover img`, chromedp.ByQuery),
		chromedp.WaitVisible(`#price_data`, chromedp.ByQuery),

		chromedp.Evaluate(`
            (() => {
                const getPrice = (id) => {
                    const el = document.querySelector("#" + id + " .price");
                    return el ? el.innerText.trim().replace("\n", "") : "N/A";
                };

                // Logic to get ONLY the set name text
                const setAnchor = document.querySelector('#product_name a');
                let extractedSetName = "";
                if (setAnchor) {
                    // Iterate through child nodes to find the text node only
                    // This ignores the <img> tag and its alt text
                    for (const node of setAnchor.childNodes) {
                        if (node.nodeType === Node.TEXT_NODE) {
                            extractedSetName += node.textContent;
                        }
                    }
                }

                const imgEl = document.querySelector('#product_details .cover img');
                
                return {
                    setName:  extractedSetName.trim(),
                    imageUrl: imgEl ? imgEl.src : "",
                    ungraded: getPrice("used_price"),
                    grade7:   getPrice("complete_price"),
                    grade8:   getPrice("new_price"),
                    grade9:   getPrice("graded_price"),
                    grade9_5:  getPrice("box_only_price"),
                    grade10:    getPrice("manual_only_price")
                };
            })()
        `, &result),
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
