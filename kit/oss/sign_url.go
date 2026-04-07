package oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

func GetSignURL(bucket_name, file_name string) (string, error) {
	client := OSSClient(bucket_name)

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		return "", err
	}
	signedURL, err := bucket.SignURL(file_name, oss.HTTPGet, 14400)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}
