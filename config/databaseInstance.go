package config

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"sync"
)

var databaseLock = &sync.Mutex{}
var databaseConnection *gorm.DB

func GetDatabaseInstance() *gorm.DB {
	//To prevent expensive lock operations
	//This means that the databaseConnection field is already populated
	if databaseConnection == nil {
		databaseLock.Lock()
		defer databaseLock.Unlock()

		//Only one goroutine can create the singleton instance.
		if databaseConnection == nil {
			log.Println("Creating Database instance")

			//DATABASE_URL=postgres://{user}:{password}@{hostname}:{port}/{database-name}?sslmode=disable
			//DATABASE_URL=postgres://postgres:Password@localhost:5432/postgres?sslmode=disable
			db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
			if err != nil {
				log.Printf("Database initialization: %s\n", err.Error())
				os.Exit(1)
			}

			databaseConnection = db
		} else {
			log.Println("Application Database instance already created.")
		}
	} else {
		//log.Println("Application Database instance already created.")
	}

	return databaseConnection
}
