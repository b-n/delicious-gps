package main

import (
  "github.com/stratoberry/go-gpsd"
  "gorm.io/driver/sqlite"
  "gorm.io/gorm"
  "os"
  "log"
  "time"
)

var (
  DebugLogger *log.Logger
  InfoLogger  *log.Logger
  ErrorLogger *log.Logger
)

func init() {
  DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
  InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
  ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
}

func check(e error) {
  if e != nil {
    ErrorLogger.Fatal(e)
    panic(e)
  }
}

type PositionData struct {
  gorm.Model
  Lat float64
  Lon float64
  Alt float64
  Velocity float64
  SatelliteCount int
  Time time.Time
}

func main() {
  db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
  check(err)

  db.AutoMigrate(&PositionData{})

  var satelliteCount = 0

  DebugLogger.Print(satelliteCount)

  tpvFilter := func(r interface {}) {
    report := r.(*gpsd.TPVReport)
    data := &PositionData{
      Lat: report.Lat,
      Lon: report.Lon,
      Alt: report.Alt,
      Velocity: report.Speed,
      SatelliteCount: 0,
      Time: report.Time,
    }
    db.Create(data)
    DebugLogger.Printf("TPV lat: %.4f lon: %.4f alt: %.4f", report.Lat, report.Lon, report.Alt)
  }

  skyFilter := func(r interface {}) {
    report := r.(*gpsd.SKYReport)
    satelliteCount = len(report.Satellites)
    DebugLogger.Printf("SKY device: %s tag: %s time: %s satelliteCount: %d", report.Device, report.Tag, report.Time, len(report.Satellites))
  }

  gps, err := gpsd.Dial(gpsd.DefaultAddress)
  check(err)

  gps.AddFilter("TPV", tpvFilter)
  gps.AddFilter("SKY", skyFilter)

  done := gps.Watch()
  <-done
}
