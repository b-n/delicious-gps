package persistence

import (
	"time"

	"gorm.io/gorm"
)

// PositionData is a struct which represents the data in the sqlite table
type PositionData struct {
	gorm.Model
	Lat            float64
	Lon            float64
	Alt            float64
	Velocity       float64
	SatelliteCount int
	Time           time.Time
}

func Open(database gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(database, &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&PositionData{})
	return db, nil
}
