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

const (
	ERRORED     uint8 = 125
	EXITING     uint8 = 126
	START_STATE uint8 = 127
)

type Options struct {
	Debug        bool
	DebugReports bool
	Database     string
}

var (
	db       *gorm.DB
	opts     Options
	appState uint8 = START_STATE
	paused   bool  = false
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

func main() {
	logging.Info("delcious-gps Started")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool)

	// Setup location listener
	locations, err := location.Listen(ctx, done)
	logging.Check(err)

	// Setup button input
	inputEvents, err := gpio.ListenInput(ctx, done)
	logging.Check(err)

	// Setup led output
	display := make(chan uint8)
	err = gpio.OpenOutput(ctx, done, display, START_STATE)
	logging.Check(err)

	// Setup state changers
	changeState, changePaused := NewStateChanger(display, &appState, &paused)

	// Handle UNIX Signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	quit := func() {
		if appState != EXITING {
			changeState(EXITING)
			close(display)
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

			changeState(gpsStateDict[location.CalculateState(v)])

			if appState < 3 || paused {
				break
			}

			logging.Debug("Storing Position Record")
			err = storePositionData(v, db)
			logging.Check(err)
		case <-inputEvents:
			changePaused(!paused)
		case <-sigs:
			quit()
			return
		}
	}
}
