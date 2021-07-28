package location

import (
	"context"

	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/stratoberry/go-gpsd"
)

// PositionData contains the location data from gpsd
type PositionData struct {
	TPVReport *gpsd.TPVReport
	SKYReport *gpsd.SKYReport
	Status    GPS_STATUS
}

var (
	initialized         bool
	tpvFilter           func(r interface{})
	skyFilter           func(r interface{})
	notificationChannel chan PositionData
	errorChannel        chan error
	lastSkyReport       *gpsd.SKYReport
	lastTpvReport       *gpsd.TPVReport
)

var satelliteCount = 0

func notify() {
	if !initialized {
		logging.Debug("GPS: dropped notify, not initialized")
		return
	}

	state := CalculateState(lastSkyReport, lastTpvReport)
	data := PositionData{
		TPVReport: lastTpvReport,
		SKYReport: lastSkyReport,
		Status:    state,
	}
	select {
	case notificationChannel <- data:
	default:
		logging.Debug("GPS: dropped notify, channel busy")
	}
}

func init() {
	tpvFilter = func(r interface{}) {
		lastTpvReport = r.(*gpsd.TPVReport)
		notify()
	}

	skyFilter = func(r interface{}) {
		lastSkyReport = r.(*gpsd.SKYReport)
	}
}

// Listen will start a listener for the gpsd service
func Listen(ctx context.Context, done chan bool) (chan PositionData, error) {
	notificationChannel = make(chan PositionData, 1)

	gps, err := gpsd.Dial(gpsd.DefaultAddress)
	if err != nil {
		return nil, err
	}

	gps.AddFilter("TPV", tpvFilter)
	gps.AddFilter("SKY", skyFilter)

	gps.Watch()

	initialized = true

	go func() {
		for {
			select {
			case <-ctx.Done():
				logging.Info("Stopping GPS")
				initialized = false
				close(notificationChannel)
				done <- true
				return
			}
		}
	}()

	logging.Info("Starting GPS")

	return notificationChannel, nil
}
