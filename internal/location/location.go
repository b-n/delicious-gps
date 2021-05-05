package location

import (
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
func Listen(c chan PositionData) (chan bool, error) {
	notificationChannel = c

	gps, err := gpsd.Dial(gpsd.DefaultAddress)
	if err != nil {
		return nil, err
	}

	gps.AddFilter("TPV", tpvFilter)
	gps.AddFilter("SKY", skyFilter)

	done := gps.Watch()
	return done, nil
}
