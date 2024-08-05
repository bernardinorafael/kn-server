package service

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/pkg/logger"
)

/*
* TODO: implement a method to remove a file from the bucket
 */

type s3Service struct {
	log    logger.Logger
	client *s3.Client
}

func NewS3Service(log logger.Logger, client *s3.Client) contract.FileManagerService {
	return &s3Service{log, client}
}

func (svc *s3Service) UploadFile(file io.Reader, key, bucket string) (location string, err error) {
	uploader := manager.NewUploader(svc.client)
	out, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String("image/webp"),
		Body:        file,
		ACL:         "public-read",
	})

	if err != nil {
		svc.log.Error("error uploading file to s3", "error", err)
		return "", errors.New("cannot upload image to bucket")
	}

	return out.Location, nil
}
