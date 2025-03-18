package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB_SOURCE string `mapstructure:"DB_SOURCE"`
	ADDRESS  string `mapstructure:"ADDRESS"`
}

func LoadConfig(path []string) (config Config, err error) {
	for _, p := range path {
		viper.AddConfigPath(p)
	}
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Config file not found %w", err))
		}
		return
	}

	err = viper.Unmarshal(&config)
	return
}
