package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
)

var (
	databaseLock       sync.Mutex
	databaseConnection *gorm.DB
)

func InitializeSQLiteDatabase() {
	databaseLock.Lock()
	defer databaseLock.Unlock()

	// Check if the database connection is already created
	if databaseConnection != nil {
		log.Println("Application Database instance already created.")
		return
	}

	log.Println("Creating Database instance")
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Printf("Database initialization: %s\n", err.Error())
		os.Exit(1)
	}

	databaseConnection = db
}

func InitializeDefaultPostgresDatabase() {
	InitializePostgresDatabase("localhost", "5432", "postgres", "postgres", "postgres", "disable")
}

func InitializePostgresDatabase(host, port, dbname, user, password, sslmode string) {
	databaseLock.Lock()
	defer databaseLock.Unlock()

	// Check if the database connection is already created
	if databaseConnection != nil {
		log.Println("Application Database instance already created.")
		return
	}

	log.Println("Creating Database instance")
	args := fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=%s", host, port, dbname, user, password, sslmode)
	db, err := gorm.Open(postgres.Open(args), &gorm.Config{})
	if err != nil {
		log.Printf("Database initialization: %s\n", err.Error())
		os.Exit(1)
	}

	databaseConnection = db
}

func GetDatabaseInstance() *gorm.DB {
	return databaseConnection
}
