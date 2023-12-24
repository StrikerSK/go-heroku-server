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

func InitializeDefaultSQLiteDatabase() {
	InitializeSQLiteDatabase("file::memory:?cache=shared")
}

func InitializeSQLiteDatabase(dsn string) {
	dialector := sqlite.Open(dsn)
	InitializeDatabase(dialector)
}

func InitializeDefaultPostgresDatabase() {
	InitializePostgresDatabase("localhost", "5432", "postgres", "postgres", "postgres", "disable")
}

func InitializePostgresDatabase(host, port, dbname, user, password, sslMode string) {
	args := fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=%s", host, port, dbname, user, password, sslMode)
	InitializeDatabase(postgres.Open(args))
}

func InitializeDatabase(dialector gorm.Dialector) {
	databaseLock.Lock()
	defer databaseLock.Unlock()

	// Check if the database connection is already created
	if databaseConnection != nil {
		log.Println("Application Database instance already created.")
		return
	}

	log.Println("Creating Database instance")
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Printf("Database initialization: %s\n", err.Error())
		os.Exit(1)
	}

	databaseConnection = db
}

func GetDatabaseInstance() *gorm.DB {
	return databaseConnection
}
