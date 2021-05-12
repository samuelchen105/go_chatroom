package common

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var database *gorm.DB

func InitDatabase() {
	//get env var
	user := os.Getenv("DB_USER")
	pwd := os.Getenv("DB_PWD")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	//open db
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", user, pwd, host, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	database = db
}

func GetDatabase() *gorm.DB {
	return database
}
