package handler

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
)

type chromedpHealthRequest struct {
	URL       string `json:"url"`
	TimeoutSec int   `json:"timeoutSec"`
}

func (h *Handler) ChromedpHealth(c echo.Context) error {
	req := chromedpHealthRequest{
		URL:       "https://example.com",
		TimeoutSec: 30,
	}

	_ = c.Bind(&req)

	if req.URL == "" {
		req.URL = "https://example.com"
	}
	if req.TimeoutSec <= 0 || req.TimeoutSec > 180 {
		req.TimeoutSec = 30
	}

	start := time.Now()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("window-size", "1920,1080"),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Duration(req.TimeoutSec)*time.Second)
	defer cancelTimeout()

	var title string
	var currentURL string

	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(req.URL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Title(&title),
		chromedp.Location(&currentURL),
	)

	if err != nil {
		return responseJSON(c, map[string]any{
			"ok":         false,
			"message":    "chromedp check failed",
			"targetUrl":  req.URL,
			"elapsedMs":  time.Since(start).Milliseconds(),
			"error":      err.Error(),
			"timeoutSec": req.TimeoutSec,
		})
	}

	return responseJSON(c, map[string]any{
		"ok":         true,
		"message":    "chromedp check passed",
		"targetUrl":  req.URL,
		"currentUrl": currentURL,
		"title":      title,
		"elapsedMs":  time.Since(start).Milliseconds(),
		"timeoutSec": req.TimeoutSec,
	})
}

