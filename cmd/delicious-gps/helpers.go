package main

import (
	"context"
	"time"

	"github.com/b-n/delicious-gps/internal/location"
	"github.com/b-n/delicious-gps/internal/mode"
	"github.com/b-n/delicious-gps/internal/persistence"
	"github.com/b-n/delicious-gps/simple_button"
)

type AppState uint8

const (
	INITIALISING AppState = iota
	WAITING_GPS
	RUNNING_AREA
	PAUSED_AREA
	RUNNING_POI
	ERRORED
	EXITING
)

var (
	stateMessage = map[AppState]string{
		INITIALISING: "Initializing",
		WAITING_GPS:  "Waiting on 3D GPS Fix",
		RUNNING_AREA: "Running (Area mode)",
		PAUSED_AREA:  "Paused (Area mode)",
		RUNNING_POI:  "Running (POI mode)",
		ERRORED:      "UNKNOWN/ERROR",
		EXITING:      "Exiting",
	}
)

func waitForInput(ctx context.Context, inputChannel chan simple_button.EventPayload, millis time.Duration) uint8 {
	timeout, cancel := context.WithTimeout(ctx, millis*time.Millisecond)
	defer cancel()

	initialButtonState := uint8(0)
	for {
		select {
		case e := <-inputChannel:
			initialButtonState |= 1 << e.Id
		case <-timeout.Done():
			return initialButtonState
		}
	}
}

func ToPositionRecord(v location.PositionData, m mode.Mode, t int) *persistence.PositionData {
	tpv := *v.TPVReport
	sky := *v.SKYReport
	return &persistence.PositionData{
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
		Mode:           int(m),
		Type:           t,
	}
}
