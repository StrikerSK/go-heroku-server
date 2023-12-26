package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go-heroku-server/config/database"
	"os"
)

type Application struct {
	Port string `mapstructure:"port"`
}

type Authorization struct {
	ExcludedPaths []string `mapstructure:"ExcludedPaths"`
	Encoding      string   `mapstructure:"Encoding"`
	Expiration    int      `mapstructure:"Expiration"`
}

type ApplicationConfiguration struct {
	Application   Application                    `mapstructure:"Application"`
	Authorization Authorization                  `mapstructure:"Authorization"`
	Database      database.DatabaseConfiguration `mapstructure:"Database"`
}

// ReadConfiguration - read file from the current directory and marshal into the conf config struct.
func ReadConfiguration() *ApplicationConfiguration {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}

	configuration := &ApplicationConfiguration{
		Application: Application{
			Port: "8080",
		},
	}

	err = viper.Unmarshal(configuration)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		os.Exit(-1)
	}

	return configuration
}
