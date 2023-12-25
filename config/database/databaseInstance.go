package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func CreateDefaultSQLiteDatabase() *gorm.DB {
	return CreateSQLiteDatabase("file::memory:?cache=shared")
}

func CreateSQLiteDatabase(dsn string) *gorm.DB {
	dialector := sqlite.Open(dsn)
	return createDatabase(dialector)
}

func CreateDefaultPostgresDatabase() *gorm.DB {
	return CreatePostgresDatabase("localhost", "5432", "postgres", "postgres", "postgres", "disable")
}

func CreatePostgresDatabase(host, port, dbname, user, password, sslMode string) *gorm.DB {
	args := fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=%s", host, port, dbname, user, password, sslMode)
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

func CreateDB(configuration PostgresDatabaseConfiguration) *gorm.DB {
	switch configuration.databaseType {
	case "sqlite":
		host := configuration.databaseHost

		if host == "" {
			panic("database host not provided")
		}

		return CreateSQLiteDatabase(host)
	default:
		panic("database type not recognized")
	}
}
