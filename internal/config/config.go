package config

import (
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	Port string `mapstructure:"PORT"`
	Mode string `mapstructure:"ENVIRONMENT"`
	Name string `mapstructure:"NAME"`

	DSN    string `mapstructure:"DB_POSTGRES_DSN"`
	DBName string `mapstructure:"DB_NAME"`
	DBUrl  string `mapstructure:"DB_URL"`

	JWTSecret           string        `mapstructure:"JWT_SECRET"`
	JwtExpiresIn        int           `mapstructure:"JWT_EXPIRES"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`

	AWSAccessKey string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion    string `mapstructure:"AWS_REGION"`
	AWSBucket    string `mapstructure:"AWS_BUCKET_NAME"`

	TwilioServiceID string `mapstructure:"TWILIO_SERVICE_SID"`
	TwilioAuthToken string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioAccountID string `mapstructure:"TWILIO_ACCOUNT_ID"`
}

func NewConfig() (*Env, error) {
	var env Env

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}
