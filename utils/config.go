package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	DbSource      string        `mapstructure:"DB_SOURCE"`
	SecretKey     string        `mapstructure:"SECRET_KEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig(path string) (config Configuration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
