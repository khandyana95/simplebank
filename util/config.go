package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver                   string        `mapstructure:"DB_DRIVER"`
	DataSource                 string        `mapstructure:"DATA_SOURCE"`
	ServerAddress              string        `mapstructure:"SERVER_ADDRESS"`
	SecretKey                  string        `mapstructure:"ACCESS_KEY"`
	TokenExpiryDuration        time.Duration `mapstructure:"TOKEN_EXPIRATION_DURATION"`
	RefreshTokenExpiryDuration time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRATION_DURATION"`
}

func LoadConfing(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
