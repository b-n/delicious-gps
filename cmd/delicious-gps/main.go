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
)

type Options struct {
	Debug        bool
	DebugReports bool
	CheckPowerSW bool
	Database     string
}

var (
	opts      Options
	appStatus AppState
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
}

func main() {
	appStatus = UpdateAppStatus(INITIALISING)

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

	_, modeData, modeDisplay := mode.Init()

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
	storage := make(chan interface{}, 30)
	err = persistence.Listen(opts.Database, storage)
	logging.Check(err)

	// Handle UNIX Signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	appStatus = UpdateAppStatus(RUNNING)

	quit := func() {
		if appStatus != EXITING {
			appStatus = UpdateAppStatus(EXITING)
			close(display)
			close(storage)
			mode.Close()
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

			//Implement GPS Status Check
			if v.SKYReport == nil {
				break
			}
			mode.HandleLocationEvent(v)

			select {
			case display <- gpio.OutputPayload{
				Index: 1,
				Blink: false,
				Color: PositionRecordToColor(v),
			}:
			default:
			}
		case d, ok := <-modeData:
			if !ok {
				break
			}
			select {
			case storage <- ToPositionRecord(
				d.Data.(location.PositionData),
				d.Mode,
				d.Type,
			):
			default:
				logging.Info("Unable to save record to database")
			}
		case d, ok := <-modeDisplay:
			if !ok {
				break
			}
			select {
			case display <- d:
			default:
			}
		case e := <-inputEvents:
			switch e.Id {
			case 0:
				quit()
				return
			case 1:
				switch e.Event {
				case simple_button.ON:
					mode.Use(mode.POI)
				case simple_button.OFF:
					mode.Use(mode.AREA)
				}
			case 2:
				mode.HandleInput(e)
			}
		case <-sigs:
			quit()
			return
		}
	}
}
