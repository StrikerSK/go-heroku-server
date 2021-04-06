package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Password"
	dbname   = "golangdb"
)

var DBConnection *gorm.DB

func InitializeDatabase() {
	//psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	//db, err := gorm.Open("postgres", psqlInfo)
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	DBConnection = db
}
