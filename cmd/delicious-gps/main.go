package main

import (
	"context"
	"flag"
	"os"

	"github.com/b-n/delicious-gps/internal/gpio"
	"github.com/b-n/delicious-gps/internal/location"
	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/b-n/delicious-gps/internal/persistence"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Options struct {
	ShowDebug bool
	Database  string
}

var (
	db        *gorm.DB
	opts      Options
	appState  uint8
	stateDict = map[uint8]string{
		0:   "Initializing",
		1:   "Waiting on SKY Report",
		2:   "Waiting on 3D Fix",
		3:   "Acquired 3D Fix, limited satellites (<=6)",
		4:   "Acquired 3D Fix, good satellites (>6)",
		255: "ERROR",
	}
)

func initOptions(args []string) Options {
	opts := Options{}

	flag.StringVar(&opts.Database, "database", "data.db", "the name of the database file to output to")
	flag.BoolVar(&opts.ShowDebug, "debug", false, "if true, output debug logging")

	flag.Parse()

	return opts
}

func init() {
	opts = initOptions(os.Args)
	logging.Init(opts.ShowDebug)

	gormdb, err := persistence.Open(sqlite.Open(opts.Database))
	logging.Check(err)

	db = gormdb
}

func storePositionData(v location.PositionData, db *gorm.DB) error {
	tpv := *v.TPVReport
	sky := *v.SKYReport
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup location listener
	locations := make(chan location.PositionData)
	gpsdDone, err := location.Listen(ctx, locations)
	logging.Check(err)

	// Setup status indicator
	controlsChannel := make(chan uint8)
	displayChannel, err := gpio.Open(ctx, controlsChannel)
	logging.Check(err)

	appState = 0
	displayChannel <- appState

	for {
		select {
		case v := <-locations:
			logging.Debugf("TPVReport: %+v", *v.TPVReport)
			if v.SKYReport != nil {
				logging.Debugf("SKYReport: %+v", *v.SKYReport)
			}

			if next := location.CalculateState(v); next != appState {
				appState = next
				displayChannel <- appState
				logging.Info(stateDict[appState])
			}

			if appState < 3 {
				break
			}

			err = storePositionData(v, db)
			logging.Debug("Stored Position Record")
			logging.Check(err)
		case <-gpsdDone:
			os.Exit(0)
		}
	}
}
