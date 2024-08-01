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

type service struct {
	log    logger.Logger
	client *s3.Client
}

func NewS3Service(client *s3.Client, log logger.Logger) contract.FileManagerService {
	return &service{log, client}
}

func (svc *service) UploadFile(file io.Reader, key, bucket string) (*manager.UploadOutput, error) {
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
		return nil, errors.New("cannot upload image to bucket")
	}

	return out, nil
}
