package Utilities

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port string
	ApiKey string
	DB DB
}

type DB struct {
	Address string
	Port string
	DBName string
	Username string
	Password string
}

func GetConfig(path string, configFileName string) (*Config) {
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic("Error Reading Config:" + err.Error())
	}

	viper.SetDefault("Port", "8080")
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic("Error Unmarshalling Config: " + err.Error())
	}

	return &config
}