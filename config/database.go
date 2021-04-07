package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

var DBConnection *gorm.DB

func InitializeDatabase() {
	//DATABASE_URL=postgres://{user}:{password}@{hostname}:{port}/{database-name}?sslmode=disable"
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	DBConnection = db
}
