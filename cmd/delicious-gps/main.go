package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/b-n/delicious-gps/internal/gpio"
	"github.com/b-n/delicious-gps/internal/location"
	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/b-n/delicious-gps/internal/persistence"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Options struct {
	Debug        bool
	DebugReports bool
	Database     string
}

var (
	db           *gorm.DB
	opts         Options
	appState     uint8
	gpsStateDict = map[location.GPS_STATUS]uint8{
		location.WAIT_SKY: 0,
		location.WAIT_FIX: 1,
		location.FIX_WEAK: 2,
		location.FIX_GOOD: 3,
	}
	stateDict = map[uint8]string{
		gpsStateDict[location.WAIT_SKY]: "Waiting on SKY Report",
		gpsStateDict[location.WAIT_FIX]: "Waiting on 3D Fix",
		gpsStateDict[location.FIX_WEAK]: "Acquired 3D Fix, limited satellites (<=6)",
		gpsStateDict[location.FIX_GOOD]: "Acquired 3D Fix, good satellites (>6)",
		125:                             "UNKNOWN/ERROR",
		126:                             "Exiting",
		127:                             "Initializing",
	}
)

func initOptions(args []string) Options {
	opts := Options{}

	flag.StringVar(&opts.Database, "database", "data.db", "the name of the database file to output to")
	flag.BoolVar(&opts.Debug, "debug", false, "if true, output debug logging")
	flag.BoolVar(&opts.DebugReports, "debug-reports", false, "Turns on debuging of raw reports (requires --debug)")

	flag.Parse()

	return opts
}

func init() {
	opts = initOptions(os.Args)
	logging.Init(opts.Debug)

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
	done := make(chan bool)

	// Setup location listener
	locations := make(chan location.PositionData)
	err := location.Listen(ctx, done, locations)
	logging.Check(err)

	// Setup button input
	controlsChannel := make(chan uint8)
	err = gpio.ListenInput(ctx, done, controlsChannel)
	logging.Check(err)

	// Setup led output
	displayChannel, err := gpio.OpenOutput(ctx, done, 127)
	logging.Check(err)

	// Handle UNIX Signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	// Set default state
	appState = 127
	logging.Info(stateDict[appState])
	displayChannel <- appState

	quit := func() {
		if appState != 126 {
			appState = 126
			logging.Info("Exiting")
			cancel()
			for i := 3; i > 0; i-- {
				<-done
			}
			return
		}
	}

	for {
		select {
		case v := <-locations:
			if opts.DebugReports {
				logging.Debugf("TPVReport: %+v", *v.TPVReport)
				if v.SKYReport != nil {
					logging.Debugf("SKYReport: %+v", *v.SKYReport)
				}
			}

			if nextState := gpsStateDict[location.CalculateState(v)]; nextState != appState && appState != 126 {
				appState = nextState
				displayChannel <- appState
				logging.Info(stateDict[appState])
			}

			if appState < 3 {
				break
			}

			err = storePositionData(v, db)
			logging.Debug("Stored Position Record")
			logging.Check(err)
		case <-controlsChannel:
			logging.Debug("BUTTONS")
		case <-sigs:
			quit()
			return
		}
	}
}
