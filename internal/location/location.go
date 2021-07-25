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
}

var (
	tpvFilter           func(r interface{})
	skyFilter           func(r interface{})
	notificationChannel chan PositionData
	errorChannel        chan error
	lastSkyReport       *gpsd.SKYReport
)

var satelliteCount = 0

func init() {
	tpvFilter = func(r interface{}) {
		report := r.(*gpsd.TPVReport)
		notificationChannel <- PositionData{
			TPVReport: report,
			SKYReport: lastSkyReport,
		}
	}

	skyFilter = func(r interface{}) {
		report := r.(*gpsd.SKYReport)
		lastSkyReport = report
	}
}

// Listen will start a listener for the gpsd service
func Listen(ctx context.Context, done chan bool, c chan PositionData) error {
	notificationChannel = c

	gps, err := gpsd.Dial(gpsd.DefaultAddress)
	if err != nil {
		return err
	}

	gps.AddFilter("TPV", tpvFilter)
	gps.AddFilter("SKY", skyFilter)

	gps.Watch()

	go func() {
		for {
			select {
			case <-ctx.Done():
				logging.Info("Stopping GPS")
				done <- true
				return
			}
		}
	}()

	logging.Info("Starting GPS")

	return nil
}
