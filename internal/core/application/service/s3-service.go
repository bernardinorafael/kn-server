package service

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	env "github.com/bernardinorafael/kn-server/internal/config"
	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
)

type service struct {
	log    *slog.Logger
	env    *env.Env
	client *s3.Client
}

func NewS3Service(client *s3.Client, env *env.Env, log *slog.Logger) contract.FileManagerService {
	return &service{log, env, client}
}

func (svc *service) UploadFile(file io.Reader, fileName string) (*manager.UploadOutput, error) {
	uploader := manager.NewUploader(svc.client)
	out, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(svc.env.AWSBucket),
		Key:         aws.String(fileName),
		ContentType: aws.String("image/*"),
		ACL:         "public-read",
		Body:        file,
	})

	if err != nil {
		return nil, errors.New("cannot upload image to bucket")
	}

	return out, nil
}
