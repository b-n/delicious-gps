package mode

import (
	"github.com/b-n/delicious-gps/internal/gpio"
	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/b-n/delicious-gps/simple_button"
)

type Mode uint8

type ModeHandler interface {
	HandleLocationEvent(e interface{})
	HandleInput(e simple_button.EventPayload) *gpio.OutputPayload
}

type LocationData struct {
	Data interface{}
	Mode Mode
	Type int
}

const (
	UNASSIGNED Mode = iota
	AREA
	POI
)

var (
	initialized bool
	am          ModeHandler
	pm          ModeHandler
	activeMode  ModeHandler
	dataChannel chan LocationData
)

func writeDataChannel(d interface{}, m Mode, t int) {
	go func(d interface{}, m Mode, t int, i *bool) {
		if *i {
			dataChannel <- LocationData{d, m, t}
		}
	}(d, m, t, &initialized)
}

func Init() (Mode, chan LocationData) {
	dataChannel = make(chan LocationData)
	am = NewAreaMode()
	pm = NewPoiMode()
	initialized = true
	return UNASSIGNED, dataChannel
}

func Use(m Mode) Mode {
	if m == AREA {
		logging.Info("Area mode active")
		activeMode = am
		return AREA
	} else if m == POI {
		logging.Info("POI mode active")
		activeMode = pm
		return POI
	}
	return UNASSIGNED
}

func Close() {
	if initialized {
		initialized = false
		close(dataChannel)
	}
}

func HandleLocationEvent(e interface{}) {
	activeMode.HandleLocationEvent(e)
}

func HandleInput(e simple_button.EventPayload) *gpio.OutputPayload {
	return activeMode.HandleInput(e)
}