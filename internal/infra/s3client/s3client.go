package s3client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	env "github.com/bernardinorafael/kn-server/internal/config"
)

func New(env *env.Env) (*s3.Client, error) {
	credential := credentials.NewStaticCredentialsProvider(env.AWSAccessKey, env.AWSSecretKey, "")
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credential),
		config.WithRegion(env.AWSRegion),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return client, nil
}
