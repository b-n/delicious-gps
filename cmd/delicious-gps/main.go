package main

import (
	"os"

	"github.com/b-n/delicious-gps/internal/config"
	"github.com/b-n/delicious-gps/internal/location"
	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/b-n/delicious-gps/internal/persistence"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	opts config.Options
)

func init() {
	opts = config.Init(os.Args)
	logging.Init(opts.ShowDebug)

	gormdb, err := persistence.Open(sqlite.Open(opts.Database))
	logging.Check(err)

	db = gormdb
}

func nextState(currentState int, positionData location.PositionData) int {
	newState := currentState

	haveSkyReport := positionData.SKYReport != nil
	have3DFix := (*positionData.TPVReport).Mode == 3

	switch {
	case (haveSkyReport && have3DFix):
		newState = 3
		logging.Info("Acquired 3D Fix, running...")
		break
	case (haveSkyReport && !have3DFix):
		newState = 2
		logging.Info("Waiting on 3D Fix")
		break
	case (!haveSkyReport):
		newState = 1
		logging.Info("Waiting on SKY Report")
		break
	}
	logging.Debugf("newState: %+v", newState)
	return newState
}

func storePositionData(pos location.PositionData, db *gorm.DB) error {
	tpv := *pos.TPVReport
	sky := *pos.SKYReport
	result := db.Create(&persistence.PositionData{
		Lon:            tpv.Lon,
		Lat:            tpv.Lat,
		Alt:            tpv.Alt,
		Velocity:       tpv.Speed,
		SatelliteCount: len(sky.Satellites),
		Time:           tpv.Time,
		ErrorLon:       tpv.Epx,
		ErrorLat:       tpv.Epy,
		ErrorAlt:       tpv.Epv,
		ErrorVelocity:  tpv.Eps,
	})

	return result.Error
}

func main() {
	logging.Info("delcious-gps Started")

	locations := make(chan location.PositionData)
	//currentStatus := make(chan status.State)

	gpsdDone, err := location.Listen(locations)
	logging.Check(err)

	state := 0

	for {
		select {
		case v := <-locations:
			logging.Debugf("TPVReport: %+v", *v.TPVReport)
			if v.SKYReport != nil {
				logging.Debugf("SKYReport: %+v", *v.SKYReport)
			}
			logging.Debugf("CurrentState: %+v", state)
			state = nextState(state, v)

			if state < 3 {
				break
			}

			err = storePositionData(v, db)
			logging.Check(err)

			logging.Debug("Processing location data")
		case <-gpsdDone:
			os.Exit(0)
		}
	}
}