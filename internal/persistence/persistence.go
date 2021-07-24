package persistence

import (
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
