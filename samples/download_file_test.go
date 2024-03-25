package main

import (
	"fmt"
	"init-container-s3/pkg/s3"
	"testing"
)

func TestBase64(t *testing.T) {
}

func TestDownloadFile(t *testing.T) {
	endpoint := "127.0.0.1:9000"
	bucket := "piston-algo-model"
	key := "rv-valuation"
	targetDir := "/dev/workspace/init-container-s3/models"
	region := "us-east-1"
	forcePathStyle := true

	err := s3.DownloadDir(endpoint, region, targetDir, bucket, key, forcePathStyle)
	if err != nil {
		fmt.Printf("Couldn't download file: %v", err)
		panic(err)
		return
	}
}
