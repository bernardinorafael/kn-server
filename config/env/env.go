package env

import "github.com/spf13/viper"

var Env *Config

type Config struct {
	Port         string `mapstructure:"PORT"`
	Mode         string `mapstructure:"ENVIRONMENT"`
	DSN          string `mapstructure:"DB_POSTGRE_DSN"`
	JwtSecret    string `mapstructure:"JWT_SECRET"`
	JwtExpiresIn int    `mapstructure:"JWT_EXPIRES"`
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
