package database

type DatabaseConfiguration struct {
	DatabaseType     string `mapstructure:"Type"`
	DatabaseHost     string `mapstructure:"Host"`
	DatabasePort     string `mapstructure:"Port"`
	DatabaseName     string `mapstructure:"Name"`
	DatabaseUsername string `mapstructure:"Username"`
	DatabasePassword string `mapstructure:"Password"`
}
