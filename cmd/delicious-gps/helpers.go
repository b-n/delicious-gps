package main

import (
	"context"
	"time"

	"github.com/b-n/delicious-gps/internal/location"
	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/b-n/delicious-gps/internal/mode"
	"github.com/b-n/delicious-gps/internal/persistence"
	"github.com/b-n/delicious-gps/simple_button"
	"github.com/lucasb-eyer/go-colorful"
)

type AppState uint8

const (
	INITIALISING AppState = iota
	RUNNING
	ERRORED
	EXITING
)

type DomainRangeFunc func(v float64) float64

func DomainToRange(d []float64, r []float64) DomainRangeFunc {
	return func(v float64) float64 {
		if v < d[0] {
			return r[0]
		}
		if v > d[1] {
			return r[1]
		}
		pos := (v - d[0]) / (d[1] - d[0])
		return (r[1]-r[0])*pos + r[0]
	}
}

var (
	stateMessage = map[AppState]string{
		INITIALISING: "Initializing",
		RUNNING:      "Running",
		ERRORED:      "UNKNOWN/ERROR",
		EXITING:      "Exiting",
	}
	gpsColor DomainRangeFunc = DomainToRange(
		[]float64{10000, 50000},
		[]float64{180, 0},
	)
)

func UpdateAppStatus(newStatus AppState) AppState {
	logging.Info(stateMessage[newStatus])
	return newStatus
}

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

func PositionRecordToColor(v location.PositionData) uint32 {
	tpv := *v.TPVReport
	totalError := tpv.Epx * tpv.Epy * tpv.Epv
	r, g, b := colorful.Hcl(gpsColor(totalError), 0.65, 0.35).Clamped().RGB255()
	return uint32(r)<<16 | uint32(g)<<8 | uint32(b)
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
