package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func CreateSQLiteDatabase(configuration DatabaseConfiguration) *gorm.DB {
	dialector := sqlite.Open(configuration.databaseHost)
	return createDatabase(dialector)
}
func CreatePostgresDatabase(configuration DatabaseConfiguration) *gorm.DB {
	host := configuration.databaseHost
	port := configuration.databasePort
	dbName := configuration.databaseName
	username := configuration.databaseUsername
	password := configuration.databasePassword
	sslMode := "disable"

	args := fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=%s", host, port, dbName, username, password, sslMode)
	return createDatabase(postgres.Open(args))
}

func createDatabase(dialector gorm.Dialector) *gorm.DB {
	log.Println("Creating Database instance")
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Printf("Database initialization: %s\n", err.Error())
		os.Exit(1)
	}

	return db
}

func CreateDB(configuration DatabaseConfiguration) *gorm.DB {
	switch configuration.databaseType {
	case "sqlite":
		if configuration.databaseHost == "" {
			panic("database host not provided")
		}

		return CreateSQLiteDatabase(configuration)
	case "postgres":
		if configuration.databaseHost == "" {
			panic("database host not provided")
		}

		if configuration.databasePort == "" {
			panic("database port not provided")
		}

		if configuration.databaseName == "" {
			panic("database name not provided")
		}

		if configuration.databaseUsername == "" {
			panic("database username not provided")
		}

		if configuration.databasePassword == "" {
			panic("database password not provided")
		}

		return CreatePostgresDatabase(configuration)
	default:
		panic("database type not recognized")
	}
}
