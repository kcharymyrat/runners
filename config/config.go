package config

import (
	"log"

	"github.com/spf13/viper"
)

func InitConfig(fileName string) *viper.Viper {
	config := viper.New()
	config.SetConfigFile(fileName)
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while parsing configuration file: %v, %v", fileName, err)
	}
	return config
}
