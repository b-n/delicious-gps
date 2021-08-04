package persistence

import (
	"context"

	"github.com/b-n/delicious-gps/internal/logging"
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

func Listen(ctx context.Context, done chan bool, db *gorm.DB, data chan interface{}) {
	go func() {
		for {
			select {
			case d := <-data:
				if result := db.Create(d); result.Error != nil {
					logging.Infof("Failed to save db record")
				}
			case <-ctx.Done():
				done <- true
				return
			}
		}
	}()
}
