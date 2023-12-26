package database

import (
	"fmt"
	"go-heroku-server/constants"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func createSQLiteDatabase(configuration DatabaseConfiguration) *gorm.DB {
	dialector := sqlite.Open(configuration.DatabaseHost)
	return createDatabase(dialector)
}
func createPostgresDatabase(configuration DatabaseConfiguration) *gorm.DB {
	host := configuration.DatabaseHost
	port := configuration.DatabasePort
	dbName := configuration.DatabaseName
	username := configuration.DatabaseUsername
	password := configuration.DatabasePassword
	sslMode := "disable"

	args := fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=%s", host, port, dbName, username, password, sslMode)
	dialector := postgres.Open(args)
	return createDatabase(dialector)
}

func createDatabase(dialector gorm.Dialector) *gorm.DB {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Printf("Database initialization: %s\n", err.Error())
		os.Exit(1)
	}

	return db
}

func CreateDB(configuration DatabaseConfiguration) *gorm.DB {
	switch configuration.DatabaseType {
	case constants.SQLiteDatabase:
		if configuration.DatabaseHost == "" {
			panic("database host not provided")
		}
		log.Println("Creating SQLite instance")
		return createSQLiteDatabase(configuration)
	case constants.PostgresDatabase:
		if configuration.DatabaseHost == "" {
			panic("database host not provided")
		}

		if configuration.DatabasePort == "" {
			panic("database port not provided")
		}

		if configuration.DatabaseName == "" {
			panic("database name not provided")
		}

		if configuration.DatabaseUsername == "" {
			panic("database username not provided")
		}

		if configuration.DatabasePassword == "" {
			panic("database password not provided")
		}

		log.Println("Creating Postgres instance")
		return createPostgresDatabase(configuration)
	default:
		panic("database type not recognized")
	}
}
