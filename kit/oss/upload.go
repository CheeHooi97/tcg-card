package oss

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func Upload(bucketName, upload_file_name string, fileByte []byte) error {
	client := OSSClient(bucketName)
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	if err := bucket.PutObject(upload_file_name, bytes.NewReader(fileByte)); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func RetrieveSignedURL(bucketName, upload_file_name string) (string, error) {
	client := OSSClient(bucketName)

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	signedURL, err := bucket.SignURL(upload_file_name, oss.HTTPGet, 24*60*60)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func ProcessImageUrl(imgUrl string, id string) ([]byte, string, error) {
	// 1. Download the image
	resp, err := http.Get(imgUrl)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// 2. Read into bytes
	fileByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// 3. Extract Filename and Extension from the URL path
	// path.Base gets the last part of the URL (e.g., "60.jpg")
	urlPath := strings.Split(imgUrl, "?")[0] // Strip query params if any
	fullFileName := path.Base(urlPath)

	imageFileExt := filepath.Ext(fullFileName)
	imageFileName := strings.TrimSuffix(fullFileName, imageFileExt)

	// 4. Validate Content Type (as a safety check)
	imageFileContentType := http.DetectContentType(fileByte)
	fmt.Printf("Detected Content-Type: %s\n", imageFileContentType)

	// 5. Create the Final Upload Name
	// Logic: "item_photo1_" + ID + "_" + OriginalName + Ext
	uploadFileName := "pokemon_" + id + "_" + imageFileName + imageFileExt

	return fileByte, uploadFileName, nil
}
