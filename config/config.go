package config

import (
	"github.com/spf13/viper"
)

func ReadConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
