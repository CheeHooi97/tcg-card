package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"pkm/utils"

	"github.com/labstack/echo/v4"
)

func pythonScraperBaseURL() string {
	base := strings.TrimSpace(os.Getenv("PY_SCRAPER_BASE_URL"))
	if base == "" {
		base = "http://127.0.0.1:8010"
	}
	return strings.TrimRight(base, "/")
}

func forwardToPythonScraper(path string, rawBody []byte) (map[string]any, error) {
	targetURL := pythonScraperBaseURL() + path

	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(rawBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload map[string]any
	if len(respBody) == 0 {
		return map[string]any{
			"ok":      false,
			"message": "empty response from python scraper",
		}, nil
	}

	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("invalid json from python scraper: %w", err)
	}

	payload["pythonStatusCode"] = resp.StatusCode
	payload["pythonEndpoint"] = targetURL
	return payload, nil
}

func (h *Handler) InsertBGSV2(c echo.Context) error {
	var req struct {
		Url []string `json:"url" validate:"required"`
	}
	if msg, err := utils.ValidateRequest(c, &req); err != nil {
		return responseValidationError(c, msg)
	}
	rawBody, _ := json.Marshal(map[string]any{"url": req.Url})

	result, err := forwardToPythonScraper("/scrape/bgs", rawBody)
	if err != nil {
		return responseJSON(c, map[string]any{
			"ok":      false,
			"message": "bgs-v2 proxy failed",
			"error":   err.Error(),
		})
	}

	return responseJSON(c, result)
}

func (h *Handler) InsertCGCV2(c echo.Context) error {
	var req struct {
		Url string `json:"url" validate:"required"`
	}
	if msg, err := utils.ValidateRequest(c, &req); err != nil {
		return responseValidationError(c, msg)
	}
	rawBody, _ := json.Marshal(map[string]any{"url": req.Url})

	result, err := forwardToPythonScraper("/scrape/cgc", rawBody)
	if err != nil {
		return responseJSON(c, map[string]any{
			"ok":      false,
			"message": "cgc-v2 proxy failed",
			"error":   err.Error(),
		})
	}

	return responseJSON(c, result)
}

func (h *Handler) InsertPSAV2(c echo.Context) error {
	var req struct {
		Urls []string `json:"urls" validate:"required"`
	}
	if msg, err := utils.ValidateRequest(c, &req); err != nil {
		return responseValidationError(c, msg)
	}
	rawBody, _ := json.Marshal(map[string]any{"urls": req.Urls})

	result, err := forwardToPythonScraper("/scrape/psa", rawBody)
	if err != nil {
		return responseJSON(c, map[string]any{
			"ok":      false,
			"message": "psa-v2 proxy failed",
			"error":   err.Error(),
		})
	}

	fmt.Println("result:", result)

	return responseJSON(c, result)
}

func (h *Handler) InsertTAGV2(c echo.Context) error {
	var req struct {
		Urls []string `json:"urls" validate:"required"`
	}
	if msg, err := utils.ValidateRequest(c, &req); err != nil {
		return responseValidationError(c, msg)
	}
	rawBody, _ := json.Marshal(map[string]any{"urls": req.Urls})

	result, err := forwardToPythonScraper("/scrape/tag", rawBody)
	if err != nil {
		return responseJSON(c, map[string]any{
			"ok":      false,
			"message": "tag-v2 proxy failed",
			"error":   err.Error(),
		})
	}

	return responseJSON(c, result)
}
