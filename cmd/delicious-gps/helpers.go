package main

import (
	"github.com/b-n/delicious-gps/internal/location"
	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/b-n/delicious-gps/internal/persistence"
	"gorm.io/gorm"
)

var (
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
		ERRORED:                         "UNKNOWN/ERROR",
		EXITING:                         "Exiting",
		START_STATE:                     "Initializing",
	}
)

type StateChanger func(nextState uint8)
type PauseChanger func(nextPaused bool)

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

func displayValue(s uint8, p bool) uint8 {
	v := s
	if p {
		v |= 128
	}
	return v
}

func NewStateChanger(displayChannel chan uint8, state *uint8, paused *bool) (StateChanger, PauseChanger) {
	logging.Info(stateDict[appState])

	return func(nextState uint8) {
			if nextState == *state {
				return
			}

			logging.Info(stateDict[nextState])

			*state = nextState
			if *state != EXITING {
				displayChannel <- displayValue(*state, *paused)
			}
		}, func(nextPaused bool) {
			if nextPaused == *paused {
				return
			}

			logging.Infof("Changing paused to %t", nextPaused)

			*paused = nextPaused
			if *state != EXITING {
				displayChannel <- displayValue(*state, *paused)
			}
		}
}
