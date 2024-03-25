package s3

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func DownloadFile(endpoint, region, targetPath, bucketName, key string, forcePathStyle bool) error {
	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	sess := session.Must(session.NewSession(aws.NewConfig().
		WithEndpoint(endpoint).
		WithRegion(region).
		WithDisableSSL(false).
		WithS3ForcePathStyle(forcePathStyle).
		WithCredentials(credentials.NewCredentials(&credentials.EnvProvider{}))))
	downloader := s3manager.NewDownloader(sess)

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, err = downloader.Download(
		f,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		},
	)
	if err == nil {
		log.Println("save file ", targetPath, " succeed")
	}
	return err
}

func DownloadDir(endpoint, region, targetPath, bucketName, key string, forcePathStyle bool) error {
	conf := aws.NewConfig().
		WithEndpoint(endpoint).
		WithDisableSSL(true).
		WithRegion(region).
		WithDisableParamValidation(true).
		WithS3UseARNRegion(true).
		WithS3ForcePathStyle(forcePathStyle).
		WithCredentialsChainVerboseErrors(true).
		WithDisableComputeChecksums(true).
		WithCredentials(credentials.NewCredentials(&credentials.EnvProvider{}))

	sess := session.Must(session.NewSession(conf))
	downloader := s3manager.NewDownloader(sess)

	svc := s3.New(sess, conf)
	params := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
		Prefix: &key,
	}

	objectsV2, err := svc.ListObjectsV2WithContext(aws.BackgroundContext(), params)
	if err != nil {
		return err
	}
	total := len(objectsV2.Contents)
	log.Println("download file size:", total)
	for i := range objectsV2.Contents {
		obj := objectsV2.Contents[i]
		if *obj.Size == 0 {
			continue
		}
		idx := i

		suffix, _ := strings.CutPrefix(*obj.Key, key)
		fn := path.Join(targetPath, suffix)
		parent := filepath.Dir(fn)
		if _, err := os.Stat(parent); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(parent, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}

		f, err := os.Create(fn)
		if err != nil {
			log.Panic(err)
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		_, err = downloader.Download(
			f,
			&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    obj.Key,
			},
		)
		if err != nil {
			log.Panic(err)
		}
		log.Println("save file ", fn, " succeed current index ", idx, " total ", total)
	}

	return err
}
