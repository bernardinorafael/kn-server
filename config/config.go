package config

import (
	"time"

	"github.com/spf13/viper"
)

var Env *EnvFile

type EnvFile struct {
	Port      string `mapstructure:"PORT"`
	Mode      string `mapstructure:"ENVIRONMENT"`
	Name      string `mapstructure:"NAME"`
	Debug     bool   `mapstructure:"DEBUG"`
	LogToFile string `mapstructure:"LOG_TO_FILE"`

	DSN    string `mapstructure:"DB_POSTGRES_DSN"`
	DBName string `mapstructure:"DB_NAME"`
	DBUrl  string `mapstructure:"DB_URL"`

	JWTSecret           string        `mapstructure:"JWT_SECRET"`
	JwtExpiresIn        int           `mapstructure:"JWT_EXPIRES"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func GetConfigEnv() (*EnvFile, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&Env)
	if err != nil {
		return nil, err
	}

	return Env, nil
}
