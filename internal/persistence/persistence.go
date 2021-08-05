package persistence

import (
	"github.com/b-n/delicious-gps/internal/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(database gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(database, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&PositionData{})
	return db, nil
}

func Listen(database string, data chan interface{}) error {
	db, err := Open(sqlite.Open(database))
	if err != nil {
		return err
	}

	go func() {
		for d := range data {
			logging.Debugf("Writing data")
			if result := db.Create(d); result.Error != nil {
				logging.Infof("Failed to save db record")
			}
		}
	}()
	return nil
}
