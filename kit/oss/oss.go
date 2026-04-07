package oss

import (
	"log"
	"pkm/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func OSSClient(bucket string) *oss.Client {
	OSSEndpoint := config.OSSEndpoint
	OSSAccessKeyID := config.OSSAccessKeyID
	OSSAccessKeySecret := config.OSSAccessKeySecret

	endpoint := OSSEndpoint
	if endpoint == "" {
		log.Fatal("Error: OSS_ENDPOINT environment variable is empty")
	}

	client, err := oss.New(OSSEndpoint, OSSAccessKeyID, OSSAccessKeySecret)
	if err != nil {
		log.Fatalf("oss.NewClient: %v", err)
	}

	// err = client.CreateBucket(bucket, oss.ACL(oss.ACLPublicRead), oss.StorageClass(oss.StorageStandard))

	// if err != nil {
	// 	// If bucket already exists, it's not a fatal error
	// 	if strings.Contains(err.Error(), "BucketAlreadyExists") {
	// 		fmt.Println("Bucket already exists, skipping creation.")
	// 	} else {
	// 		fmt.Printf("Error creating bucket: %v\n", err)
	// 	}
	// } else {
	// 	fmt.Printf("Successfully created bucket: %s\n", config.OSSBucket)
	// }
	return client
}
