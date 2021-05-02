package main

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

var db *gorm.DB

func check(e error) {
	if e != nil {
		errorLogger.Fatal(e)
		panic(e)
	}
}

func init() {
	timeFormat := log.Ldate | log.Ltime

	gormdb, err := persistence.init(sqlite.Open("data.db"))
	check(err)

	db = gormdb

	debugLogger = log.New(os.Stdout, "DEBUG: ", timeFormat)
	infoLogger = log.New(os.Stdout, "INFO: ", timeFormat)
	errorLogger = log.New(os.Stdout, "ERROR: ", timeFormat)
}

func main() {
	_, err := location.listen()

	debugLogger.Printf("here")
}
