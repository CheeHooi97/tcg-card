package handler

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"pkm/config"
	"pkm/config/set"
	"pkm/errcode"
	"pkm/kit/oss"
	"pkm/model"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateSet(c echo.Context) error {
	if err := c.Request().ParseMultipartForm(32 << 20); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid multipart form")
	}

	var fileHeaders []*multipart.FileHeader
	if c.Request().MultipartForm != nil {
		fileHeaders = c.Request().MultipartForm.File["photo"]
	}

	setLists := set.SetName

	for _, fileHeader := range fileHeaders {
		for setList := range setLists {
			set := model.NewSet()
			var fileName string

			fmt.Println("fileHeader.Filename: ", fileHeader.Filename)
			if !strings.Contains(fileHeader.Filename, setList) {
				continue
			}

			file, err := fileHeader.Open()
			if err != nil {
				log.Printf("Failed to open uploaded file: %v", err)
				continue
			}

			defer file.Close()

			fileByte, err := io.ReadAll(file)
			if err != nil {
				log.Printf("Failed to read uploaded file: %v", err)
				continue
			}

			ext := filepath.Ext(fileHeader.Filename)

			fileName = "pokemon_set_" + setList + ext

			if err := oss.Upload(config.OSSBucket, fileName, fileByte); err != nil {
				return responseError(c, errcode.InternalServerError)
			}

			set.Name = setList
			set.PhotoUrl = fileName

			if err := h.Set.Create(set); err != nil {
				return responseError(c, errcode.InternalServerError)
			}

			break
		}
	}

	return responseJSON(c, echo.Map{
		"lists": true,
	})
}
