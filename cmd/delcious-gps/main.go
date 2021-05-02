package main

import (
	"log"
	"os"

	"github.com/b-n/delicious-gps/internal/location"
	"github.com/b-n/delicious-gps/internal/persistence"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	infoLogger  *log.Logger
	debugLogger *log.Logger
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

	gormdb, err := persistence.Open(sqlite.Open("data.db"))
	check(err)

	db = gormdb

	debugLogger = log.New(os.Stdout, "DEBUG: ", timeFormat)
	infoLogger = log.New(os.Stdout, "INFO: ", timeFormat)
	errorLogger = log.New(os.Stdout, "ERROR: ", timeFormat)
}

func main() {
	infoLogger.Printf("delcious-gps Started")
	locations := make(chan location.PositionData)

	gpsdDone, err := location.Listen(locations)
	check(err)

	for {
		select {
		case v := <-locations:
			debugLogger.Printf("Location lon: %.4f lat: %.4f alt: %.4f", v.Lon, v.Lat, v.Alt)
		case <-gpsdDone:
			os.Exit(0)
		}
	}

	debugLogger.Printf("here")
}
