package config

import (
	"github.com/spf13/viper"
	"log"
)

func InitConfig(fileName string) *viper.Viper {
	config := viper.New()
	config.SetConfigFile(fileName)
	config.AddConfigPath(".")
	config.AddConfigPath("$HOME")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading config file, ", err)
	}
	return config
}
