package service

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/pkg/logger"
)

type sesService struct {
	log    logger.Logger
	client *s3.Client
}

func NewSESService(log logger.Logger, client *s3.Client) contract.EmailNotifier {
	return &sesService{log, client}
}

func (svc sesService) Notify(to []string, subject string, body string) error {
	// TODO implement me
	panic("implement me")
}
