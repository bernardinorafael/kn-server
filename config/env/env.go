package env

import "github.com/spf13/viper"

var Env *Config

type Config struct {
	Mode string `mapstructure:"ENVIRONMENT"`
	Port string `mapstructure:"PORT"`

	DSN string `mapstructure:"DB_POSTGRE_DSN"`
}

func LoadConfig() (*Config, error) {
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
