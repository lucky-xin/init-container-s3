package main

import (
	"fmt"
	"init-container-s3/pkg/s3"
	"os"
)

func Env(name, defaultVal string) string {
	val := os.Getenv(name)
	if val != "" {
		return val
	}
	return defaultVal
}

func main() {
	endpoint := Env("S3_ENDPOINT", "127.0.0.1:9000")
	bucket := Env("S3_BUCKET", "piston-algo-model")
	key := Env("S3_KEY", "rv-valuation")
	targetDir := Env("S3_TARGET", "/tmp/s3")
	ft := Env("S3_TYPE", "dir")
	region := Env("S3_REGION", "us-east-1")
	forcePathStyle := Env("S3_FORCE_PATH_STYLE", "false") == "true"
	if "file" == ft {
		err := s3.DownloadFile(endpoint, region, targetDir, bucket, key, forcePathStyle)
		if err != nil {
			fmt.Printf("Couldn't download file: %v", err)
			panic(err)
			return
		}
	} else if "dir" == ft {
		err := s3.DownloadDir(endpoint, region, targetDir, bucket, key, forcePathStyle)
		if err != nil {
			fmt.Printf("Couldn't download file: %v", err)
			panic(err)
			return
		}
	}
}
