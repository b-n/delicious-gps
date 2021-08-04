package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/b-n/delicious-gps/internal/gpio"
	"github.com/b-n/delicious-gps/internal/location"
	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/b-n/delicious-gps/internal/mode"
	"github.com/b-n/delicious-gps/internal/persistence"
	"github.com/b-n/delicious-gps/simple_button"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Options struct {
	Debug        bool
	DebugReports bool
	CheckPowerSW bool
	Database     string
}

var (
	db        *gorm.DB
	opts      Options
	appStatus AppState = INITIALISING
	gpsState  location.GPSState
	paused    bool = false
)

func initOptions(args []string) Options {
	opts := Options{}

	flag.StringVar(&opts.Database, "database", "data.db", "the name of the database file to output to")
	flag.BoolVar(&opts.Debug, "debug", false, "if true, output debug logging")
	flag.BoolVar(&opts.DebugReports, "debug-reports", false, "Turns on debuging of raw reports (requires --debug)")
	flag.BoolVar(&opts.CheckPowerSW, "check-power-switch", true, "Check power switch on start up (and die if not enabled)")

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

	// create channels
	locations := location.Init()
	display := make(chan gpio.OutputPayload)
	inputEvents := gpio.Init(display)

	// Start GPIO input
	err := gpio.ListenInput(ctx, done)
	logging.Check(err)

	// Check power switch (if required), and quick if needed
	initialState := waitForInput(ctx, inputEvents, 500)
	if opts.CheckPowerSW && initialState&1 != 1 {
		logging.Check(errors.New("Exiting: On switch is currently off"))
	}

	_, modeData := mode.Init()

	if initialState&2 == 2 {
		mode.Use(mode.POI)
	} else {
		mode.Use(mode.AREA)
	}

	// Open the other listeners
	err = location.Listen(ctx, done)
	logging.Check(err)
	err = gpio.OpenOutput(ctx, done)
	logging.Check(err)
	storage := make(chan interface{}, 10)
	persistence.Listen(ctx, done, db, storage)

	// Handle UNIX Signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	quit := func() {
		if appStatus != EXITING {
			appStatus = EXITING
			close(display)
			close(storage)
			mode.Close()
			cancel()
			for i := 4; i > 0; i-- {
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

			mode.HandleLocationEvent(v)
		case d := <-modeData:
			select {
			case storage <- ToPositionRecord(
				d.Data.(location.PositionData),
				d.Mode,
				d.Type,
			):
			default:
				logging.Info("Unable to save record to database")
			}
		case e := <-inputEvents:
			switch e.Id {
			case 0:
				quit()
			case 1:
				switch e.Event {
				case simple_button.ON:
					mode.Use(mode.POI)
				case simple_button.OFF:
					mode.Use(mode.AREA)
				}
			case 2:
				if de := mode.HandleInput(e); de != nil {
					select {
					case display <- *de:
					default:
					}
				}
			}
		case <-sigs:
			quit()
			return
		}
	}
}
