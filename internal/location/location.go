package location

import (
	"context"
	"time"

	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/stratoberry/go-gpsd"
)

// PositionData contains the location data from gpsd
type PositionData struct {
	TPVReport  *gpsd.TPVReport
	SKYReport  *gpsd.SKYReport
	Status     GPSState
	CreatedAt  int64
	TotalError float64
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

func errorFromTPV(r *gpsd.TPVReport) float64 {
	return r.Epx * r.Epy * r.Epv
}

func notify() {
	if !initialized {
		logging.Debug("GPS: dropped notify, not initialized")
		return
	}

	state := CalculateState(lastSkyReport, lastTpvReport)
	data := PositionData{
		TPVReport:  lastTpvReport,
		SKYReport:  lastSkyReport,
		Status:     state,
		CreatedAt:  time.Now().Unix(),
		TotalError: errorFromTPV(lastTpvReport),
	}
	select {
	case notificationChannel <- data:
	default:
		logging.Debug("GPS: dropped notify, channel busy")
	}
}

func Init() chan PositionData {
	initialized = false
	notificationChannel = make(chan PositionData, 1)

	tpvFilter = func(r interface{}) {
		lastTpvReport = r.(*gpsd.TPVReport)
		notify()
	}

	skyFilter = func(r interface{}) {
		lastSkyReport = r.(*gpsd.SKYReport)
	}

	return notificationChannel
}

// Listen will start a listener for the gpsd service
func Listen(ctx context.Context, done chan bool) error {
	gps, err := gpsd.Dial(gpsd.DefaultAddress)
	if err != nil {
		return err
	}

	gps.AddFilter("TPV", tpvFilter)
	gps.AddFilter("SKY", skyFilter)
	gps.Watch()

	initialized = true

	logging.Debug("Watching GPS")

	go func() {
		<-ctx.Done()
		logging.Debug("Stopping GPS")
		initialized = false
		close(notificationChannel)
		done <- true
	}()

	return nil
}
