package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	DBName   string `mapstructure:"SQLITE_DB"`
	DBUser   string `mapstructure:"SQLITE_USER"`
	DBPass   string `mapstructure:"SQLITE_PASSWORD"`
	Protocol string `mapstructure:"PROTOCOL"`
}

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to load config file")
		return err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to load config file")
		return err
	}

	AppConfig = &cfg
	return nil
}
