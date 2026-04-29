package handler

import (
	"context"
	"fmt"
	"net/url"
	"pkm/config/rarity"
	"pkm/model"
	"pkm/utils"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
)

const defaultBGSSportID = "477173"
const bgsMaxRetriesPerLink = 3

func (h *Handler) InsertBGS(c echo.Context) error {
	var i struct {
		Url []string `json:"url" validate:"required"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

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

	for z := 0; z < len(i.Url); z++ {
		url := i.Url[z]
		retryFromLink := 0
		retryCounts := map[int]int{}

		for {
			sessionCtx, sessionCancel := chromedp.NewContext(allocCtx)
			batchCtx, batchCancel := context.WithTimeout(sessionCtx, 30*time.Minute)

			links, sportID, err := BGSScrapLink(batchCtx, url)
			if err != nil {
				batchCancel()
				sessionCancel()
				fmt.Println("BGSScrapLink error:", err)
				break
			}

			if sportID == "" {
				sportID = defaultBGSSportID
			}

			if len(links) == 0 {
				batchCancel()
				sessionCancel()
				fmt.Println("No links found for:", url)
				break
			}

			if retryFromLink >= len(links) {
				batchCancel()
				sessionCancel()
				fmt.Printf("Retry index out of range (%d/%d), skipping URL: %s\n", retryFromLink+1, len(links), url)
				break
			}

			time.Sleep(5 * time.Second)

			fmt.Println("Link: ", links)

			shouldRetry := false
			for y := retryFromLink; y < len(links); y++ {
				currentLink := links[y]
				fmt.Printf("Scraping set link %d/%d: %s\n", y+1, len(links), currentLink.URL)

				// if y < 179 {
				// 	continue
				// }

				bgsCards, err := ScrapBGS(batchCtx, sportID, currentLink.URL)
				if err != nil {
					fmt.Printf("ScrapBGS error (link %d/%d): %v\n", y+1, len(links), err)
					retryCounts[y]++
					if retryCounts[y] <= bgsMaxRetriesPerLink {
						retryFromLink = y
						shouldRetry = true
						fmt.Printf("Retrying link %d/%d (attempt %d/%d)\n", y+1, len(links), retryCounts[y], bgsMaxRetriesPerLink)
					} else {
						fmt.Printf("Skipping failed link %d/%d after %d retries\n", y+1, len(links), bgsMaxRetriesPerLink)
						retryFromLink = y + 1
						shouldRetry = retryFromLink < len(links)
					}
					break
				}

				for x, bgsCard := range bgsCards {
					if x == 0 {
						fmt.Printf("%d of %d Lists \n", z+1, len(i.Url))
						fmt.Printf("%d of %d Sets \n", y+1, len(links))
					}

					fmt.Println("bgsCard: ", bgsCard)

					fmt.Printf("%d of %d Cards \n", x+1, len(bgsCards))

					check := h.BGS.GetByCardNameAndCardNumberAndSetId(bgsCard.CardName, bgsCard.CardNumber, bgsCard.SetID)

					if !check {
						bgs := model.NewBGS()
						bgs.CardName = bgsCard.CardName
						bgs.CardNumber = bgsCard.CardNumber
						bgs.SetId = bgsCard.SetID
						bgs.Total = bgsCard.TotalCount
						bgs.Grade1 = bgsCard.GradeCounts["1"]
						bgs.Grade1_5 = bgsCard.GradeCounts["1.5"]
						bgs.Grade2 = bgsCard.GradeCounts["2"]
						bgs.Grade2_5 = bgsCard.GradeCounts["2.5"]
						bgs.Grade3 = bgsCard.GradeCounts["3"]
						bgs.Grade3_5 = bgsCard.GradeCounts["3.5"]
						bgs.Grade4 = bgsCard.GradeCounts["4"]
						bgs.Grade4_5 = bgsCard.GradeCounts["4.5"]
						bgs.Grade5 = bgsCard.GradeCounts["5"]
						bgs.Grade5_5 = bgsCard.GradeCounts["5.5"]
						bgs.Grade6 = bgsCard.GradeCounts["6"]
						bgs.Grade6_5 = bgsCard.GradeCounts["6.5"]
						bgs.Grade7 = bgsCard.GradeCounts["7"]
						bgs.Grade7_5 = bgsCard.GradeCounts["7.5"]
						bgs.Grade8 = bgsCard.GradeCounts["8"]
						bgs.Grade8_5 = bgsCard.GradeCounts["8.5"]
						bgs.Grade9 = bgsCard.GradeCounts["9"]
						bgs.Grade9_5 = bgsCard.GradeCounts["9.5"]
						bgs.Grade10P = bgsCard.GradeCounts["10P"]
						bgs.Grade10BL = bgsCard.GradeCounts["10BL"]

						title := strings.TrimSpace(bgsCard.CardName)

						for key := range rarity.BGSRarities {
							// check if title ends with the rarity key
							if strings.HasSuffix(title, key) {
								bgs.CardName = strings.TrimSpace(title[:len(title)-len(key)])
								bgs.Rarity = key
							}
						}

						bgs.SetName = bgsCard.SetTitle

						if err := h.BGS.Create(bgs); err != nil {
							fmt.Println("BGS create error:", err)
							continue
						}

						logging := model.NewBGSLogging()
						logging.BGSId = bgs.Id
						logging.CardName = bgs.CardName
						logging.CardNumber = bgs.CardNumber
						logging.SetId = bgs.SetId
						logging.SetName = bgs.SetName
						logging.Rarity = bgs.Rarity
						logging.Description = bgs.Description
						logging.SetNumber = bgs.SetNumber
						logging.Total = bgs.Total
						logging.Grade1 = bgs.Grade1
						logging.Grade1_5 = bgs.Grade1_5
						logging.Grade2 = bgs.Grade2
						logging.Grade2_5 = bgs.Grade2_5
						logging.Grade3 = bgs.Grade3
						logging.Grade3_5 = bgs.Grade3_5
						logging.Grade4 = bgs.Grade4
						logging.Grade4_5 = bgs.Grade4_5
						logging.Grade5 = bgs.Grade5
						logging.Grade5_5 = bgs.Grade5_5
						logging.Grade6 = bgs.Grade6
						logging.Grade6_5 = bgs.Grade6_5
						logging.Grade7 = bgs.Grade7
						logging.Grade7_5 = bgs.Grade7_5
						logging.Grade8 = bgs.Grade8
						logging.Grade8_5 = bgs.Grade8_5
						logging.Grade9 = bgs.Grade9
						logging.Grade9_5 = bgs.Grade9_5
						logging.Grade10P = bgs.Grade10P
						logging.Grade10BL = bgs.Grade10BL

						if err := h.BGSLogging.Create(logging); err != nil {
							fmt.Println("BGS logging create error:", err)
							continue
						}

						// urlCheck := h.BGSUrl.GetByPath(link.URL)

						// if !urlCheck {
						// 	bgsUrl := model.NewBGSUrl()
						// 	bgsUrl.Url = link.URL

						// 	if err := h.BGSUrl.Create(bgsUrl); err != nil {
						// 		return responseError(c, errcode.InternalServerError)
						// 	}
						// }
					} else {
						dbBGS, err := h.BGS.GetDetailByCardNameAndCardNumberAndSetId(bgsCard.CardName, bgsCard.CardNumber, bgsCard.SetID)
						if err != nil {
							fmt.Println("BGS detail lookup error:", err)
							continue
						}

						logging := model.NewBGSLogging()
						if dbBGS != nil && dbBGS.Id != "" {
							logging.BGSId = dbBGS.Id
						}

						logging.CardName = bgsCard.CardName
						logging.CardNumber = bgsCard.CardNumber
						logging.SetId = bgsCard.SetID
						logging.SetName = bgsCard.SetTitle
						logging.Total = bgsCard.TotalCount
						logging.Grade1 = bgsCard.GradeCounts["1"]
						logging.Grade1_5 = bgsCard.GradeCounts["1.5"]
						logging.Grade2 = bgsCard.GradeCounts["2"]
						logging.Grade2_5 = bgsCard.GradeCounts["2.5"]
						logging.Grade3 = bgsCard.GradeCounts["3"]
						logging.Grade3_5 = bgsCard.GradeCounts["3.5"]
						logging.Grade4 = bgsCard.GradeCounts["4"]
						logging.Grade4_5 = bgsCard.GradeCounts["4.5"]
						logging.Grade5 = bgsCard.GradeCounts["5"]
						logging.Grade5_5 = bgsCard.GradeCounts["5.5"]
						logging.Grade6 = bgsCard.GradeCounts["6"]
						logging.Grade6_5 = bgsCard.GradeCounts["6.5"]
						logging.Grade7 = bgsCard.GradeCounts["7"]
						logging.Grade7_5 = bgsCard.GradeCounts["7.5"]
						logging.Grade8 = bgsCard.GradeCounts["8"]
						logging.Grade8_5 = bgsCard.GradeCounts["8.5"]
						logging.Grade9 = bgsCard.GradeCounts["9"]
						logging.Grade9_5 = bgsCard.GradeCounts["9.5"]
						logging.Grade10P = bgsCard.GradeCounts["10P"]
						logging.Grade10BL = bgsCard.GradeCounts["10BL"]

						if dbBGS != nil && dbBGS.Id != "" {
							if dbBGS.SetName != "" {
								logging.SetName = dbBGS.SetName
							}
							if dbBGS.SetNumber != "" {
								logging.SetNumber = dbBGS.SetNumber
							}
							if dbBGS.Rarity != "" {
								logging.Rarity = dbBGS.Rarity
							}
							if dbBGS.Description != "" {
								logging.Description = dbBGS.Description
							}
						}

						if err := h.BGSLogging.Create(logging); err != nil {
							fmt.Println("BGS logging create error:", err)
							continue
						}
					}
				}
				time.Sleep(5 * time.Second)
			}

			batchCancel()
			sessionCancel()

			if shouldRetry {
				fmt.Printf("Refreshing context and retrying from error link %d/%d...\n", retryFromLink+1, len(links))
				time.Sleep(2 * time.Second)
				continue
			}

			break
		}

		time.Sleep(5 * time.Second)
	}

	fmt.Println("--------------------------")
	fmt.Println("End Scrap BGS")

	return responseJSON(c, true)
}

func ScrapBGS(ctx context.Context, sportID, rawURL string) ([]BGSPopData, error) {
	var htmlContent string
	targetURL := addSportIDToBGSURL(rawURL, sportID)

	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.Sleep(2*time.Second),

		// Only bridge list -> table. Do not trigger search again,
		// otherwise Beckett navigates back to pop-report.
		chromedp.ActionFunc(func(ctx context.Context) error {
			var listCount int
			if err := chromedp.Evaluate(`document.querySelectorAll(".pop_search_list ul li a").length`, &listCount).Do(ctx); err != nil {
				return err
			}
			if listCount > 0 {
				fmt.Println("Detected list view, clicking into table...")
				return chromedp.Click(`.pop_search_list ul li a`, chromedp.ByQuery).Do(ctx)
			}
			return nil
		}),

		// 3. Now wait for the data table to appear
		chromedp.WaitVisible(`tr.rows`, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),

		// 4. Infinite Scroll Logic
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 15; i++ {
				var rowCount int
				chromedp.Evaluate(`document.querySelectorAll("tr.rows").length`, &rowCount).Do(ctx)
				if rowCount == 0 {
					break
				}

				fmt.Printf("Scrolling... %d rows\n", rowCount)
				scrollScript := fmt.Sprintf(`document.querySelectorAll("tr.rows")[%d].scrollIntoView()`, rowCount-1)
				chromedp.Evaluate(scrollScript, nil).Do(ctx)

				success := false
				for retry := 0; retry < 5; retry++ {
					time.Sleep(2 * time.Second)
					var newCount int
					chromedp.Evaluate(`document.querySelectorAll("tr.rows").length`, &newCount).Do(ctx)
					if newCount > rowCount {
						success = true
						break
					}
				}
				if !success {
					break
				}
			}
			return nil
		}),
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)

	if err != nil {
		return nil, err
	}

	return parseBGSResults(htmlContent)
}

func parseBGSResults(htmlContent string) ([]BGSPopData, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var results []BGSPopData

	// Find every row in the population table
	doc.Find("tr.rows").Each(func(i int, s *goquery.Selection) {
		// 1. Extract basic hidden input values
		setTitle := s.Find("input.set_title").AttrOr("value", "")
		totalVal := s.Find("input.card_total_value").AttrOr("value", "0")

		// 2. Extract Card Identifiers
		cardName := strings.TrimSpace(s.Find("td.test").Eq(0).Text())
		cardNumber := strings.TrimSpace(s.Find("td.test").Eq(1).Text())

		// 3. Extract SetID from the link if it exists
		link, exists := s.Find("a").Attr("href")
		setID := ""
		if exists {
			u, _ := url.Parse(link)
			setID = u.Query().Get("set_id")
		}

		// 4. Map the grade counts (10BL, 10P, 9.5, etc.)
		gradeCounts := make(map[string]string)
		s.Find("td").Each(func(j int, td *goquery.Selection) {
			gradeValue := td.Find("input.header_grade").AttrOr("value", "")
			if gradeValue == "" {
				return
			}

			// Determine the count (BGS puts counts inside <b><a href...> if they aren't zero)
			count := strings.TrimSpace(td.Find("b.popCard a").Text())
			if count == "" {
				count = strings.TrimSpace(td.Text())
			}
			// Clean up dashes or empty strings
			if count == "-" || count == "" {
				count = "0"
			}

			if gradeValue == "10" {
				headerType := td.Find("input.header_type").AttrOr("value", "")
				// Check for Black Label vs Pristine
				if strings.Contains(headerType, "Black Label") {
					gradeCounts["10BL"] = count
				} else {
					gradeCounts["10P"] = count
				}
			} else {
				gradeCounts[gradeValue] = count
			}
		})

		// 5. Build the final data object
		results = append(results, BGSPopData{
			CardName:    cardName,
			CardNumber:  cardNumber,
			TotalCount:  totalVal,
			SetTitle:    setTitle,
			SetID:       setID,
			GradeCounts: gradeCounts,
		})
	})

	return results, nil
}

type BGSUrls struct {
	URL string `json:"url"`
}

func BGSScrapLink(ctx context.Context, setName string) ([]BGSUrls, string, error) {
	var urls []BGSUrls
	var sportID string

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.beckett.com/grading/pop-report"),
		chromedp.WaitVisible(`#sport_id`, chromedp.ByID),
		chromedp.Evaluate(`document.querySelector('#sport_id')?.value || ""`, &sportID),

		// 1. Set the Sport
		chromedp.ActionFunc(func(ctx context.Context) error {
			if sportID == "" {
				sportID = defaultBGSSportID
			}
			return chromedp.SetValue(`#sport_id`, sportID, chromedp.ByID).Do(ctx)
		}),

		// 2. Clear and Type Set Name
		chromedp.Click(`#set_name`, chromedp.ByID),
		chromedp.SendKeys(`#set_name`, setName, chromedp.ByID),

		// 3. Trigger Search via ENTER key (Often more reliable than Click)
		chromedp.KeyEvent("\r"),

		// 4. WAIT for the specific links to exist in the DOM
		// We use the attribute starts with selector for 'set_match'
		chromedp.WaitVisible(`a[href*="/set_match/"]`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),

		// 5. EXTRACT using a very broad but safe selector
		chromedp.Evaluate(`
			(() => {
				// Find all anchors that link to a set_match page
				const anchors = Array.from(document.querySelectorAll('a[href*="/set_match/"]'));
				
				return anchors.map(a => {
					return {
						url: a.href,
						setName: a.innerText.trim()
					};
				});
			})()
		`, &urls),
	)

	if err != nil {
		return nil, "", fmt.Errorf("BGS Scrape failed: %v", err)
	}
	if sportID == "" {
		sportID = defaultBGSSportID
	}

	return urls, sportID, nil
}

func addSportIDToBGSURL(rawURL, sportID string) string {
	if sportID == "" {
		sportID = defaultBGSSportID
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	q := u.Query()
	q.Set("sport_id", sportID)
	u.RawQuery = q.Encode()
	return u.String()
}

type BGSPopData struct {
	CardName    string            `json:"cardName"`
	CardNumber  string            `json:"cardNumber"`
	TotalCount  string            `json:"totalCount"`
	SetTitle    string            `json:"setTitle"`
	SetID       string            `json:"setID"`
	GradeCounts map[string]string `json:"gradeCounts"`
}
