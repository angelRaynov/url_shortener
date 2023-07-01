package config

import (
	"github.com/spf13/viper"
	"log"
)

type Application struct {
	AppENV  string `mapstructure:"app_env"`
	AppMode string `mapstructure:"app_mode"`
	AppPort string `mapstructure:"app_port"`
	AppURL  string `mapstructure:"app_url"`

	JWTSecret string `mapstructure:"jwt_secret"`

	DBHost     string `mapstructure:"db_host"`
	DBKeyspace string `mapstructure:"db_keyspace"`

	RedisHost string `mapstructure:"redis_host"`
}

func New() *Application {
	viper.SetConfigFile(".config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var config Application
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}

	//TODO: remove before deploy
	log.Printf("%+v", config)
	return &config
}
