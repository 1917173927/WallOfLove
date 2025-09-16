package config

import (
	"log"

	"github.com/spf13/viper"
)

var Config = viper.New()

func init() {
	Config.AddConfigPath("./conf/config")
	Config.SetConfigName("config")
	Config.SetConfigType("yaml")
	err := Config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}