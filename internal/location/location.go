package location

import (
	"time"

	"github.com/stratoberry/go-gpsd"
)

// PositionData contains the location data from gpsd
type PositionData struct {
	Lat            float64
	Lon            float64
	Alt            float64
	Velocity       float64
	SatelliteCount int
	Time           time.Time
}

var (
	tpvFilter           func(r interface{})
	skyFilter           func(r interface{})
	notificationChannel chan PositionData
)

var satelliteCount = 0

func init() {
	tpvFilter = func(r interface{}) {
		report := r.(*gpsd.TPVReport)
		notificationChannel <- PositionData{
			Lon:            report.Lon,
			Lat:            report.Lat,
			Alt:            report.Alt,
			SatelliteCount: satelliteCount,
			Time:           report.Time,
		}
	}

	skyFilter = func(r interface{}) {
		report := r.(*gpsd.SKYReport)
		satelliteCount = len(report.Satellites)
		//log.Printf("SKY device: %s tag: %s time: %s satelliteCount: %d", report.Device, report.Tag, report.Time, len(report.Satellites))
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
