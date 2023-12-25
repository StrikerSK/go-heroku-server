package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go-heroku-server/config/database"
)

type Application struct {
	Port string `mapstructure:"port"`
}

type ApplicationConfiguration struct {
	Application Application                    `mapstructure:"Application"`
	Database    database.DatabaseConfiguration `mapstructure:"Database"`
}

// ReadConfiguration - read file from the current directory and marshal into the conf config struct.
func ReadConfiguration() *ApplicationConfiguration {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}

	configuration := &ApplicationConfiguration{
		Application: Application{
			Port: "8080",
		},
	}

	err = viper.Unmarshal(configuration)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return configuration
}
