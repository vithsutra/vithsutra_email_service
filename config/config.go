// boiler  plate code
package config

import (
	"log"

	"github.com/spf13/viper"
)

var Config struct {
	Kafka struct {
		Brokers []string
		GroupID string
		Topics  []string
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		From     string
	}
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	err := viper.Unmarshal(&Config)
	if err != nil {
		return err
	}

	log.Println("Configuration loaded successfully!")
	return nil
}
