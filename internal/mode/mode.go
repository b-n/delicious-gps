package mode

import (
	"github.com/b-n/delicious-gps/internal/gpio"
	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/b-n/delicious-gps/simple_button"
)

type Mode uint8

type ModeHandler interface {
	HandleLocationEvent(e interface{})
	HandleInput(e simple_button.EventPayload)
	Activate()
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
	initialized    bool
	am             ModeHandler
	pm             ModeHandler
	activeMode     ModeHandler
	dataChannel    chan LocationData
	displayChannel chan gpio.OutputPayload
)

func writeDataChannel(d interface{}, m Mode, t int) {
	go func(d interface{}, m Mode, t int, i *bool) {
		if *i {
			dataChannel <- LocationData{d, m, t}
		}
	}(d, m, t, &initialized)
}

func writeDisplayChannel(d gpio.OutputPayload) {
	go func(d gpio.OutputPayload, i *bool) {
		if *i {
			displayChannel <- d
		}
	}(d, &initialized)
}

func Init() (Mode, chan LocationData, chan gpio.OutputPayload) {
	dataChannel = make(chan LocationData)
	displayChannel = make(chan gpio.OutputPayload)
	am = NewAreaMode()
	pm = NewPoiMode()
	initialized = true
	return UNASSIGNED, dataChannel, displayChannel
}

func Use(m Mode) Mode {
	if m == UNASSIGNED {
		activeMode = nil
		return m
	}
	if m == AREA {
		logging.Info("Area mode active")
		activeMode = am
	} else if m == POI {
		logging.Info("POI mode active")
		activeMode = pm
	}
	activeMode.Activate()
	return m
}

func Close() {
	if initialized {
		initialized = false
		close(dataChannel)
		close(displayChannel)
	}
}

func HandleLocationEvent(e interface{}) {
	activeMode.HandleLocationEvent(e)
}

func HandleInput(e simple_button.EventPayload) {
	activeMode.HandleInput(e)
}
