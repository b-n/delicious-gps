package mode

import (
	"container/ring"

	"github.com/b-n/delicious-gps/internal/gpio"
	"github.com/b-n/delicious-gps/simple_button"
)

type AreaType struct {
	Name  string
	Color uint32
}

type AreaMode struct {
	paused bool
	areas  *ring.Ring
}

func (a *AreaMode) RegisterType(name string, color uint32) {
	j := ring.New(1)
	j.Value = AreaType{Name: name, Color: color}

	if a.areas == nil {
		a.areas = j
	} else {
		a.areas = a.areas.Link(j)
	}
}

func (a *AreaMode) HandleLocationEvent(e interface{}) {
	if !a.paused {
		writeDataChannel(e, AREA)
	}
}

func (a *AreaMode) HandleInput(e simple_button.EventPayload) *gpio.OutputPayload {
	switch e.Event {
	case simple_button.DBL_CLICK:
		a.paused = !a.paused
		return &gpio.OutputPayload{a.paused, a.areas.Value.(AreaType).Color}
	case simple_button.CLICK:
		a.areas = a.areas.Next()
		return &gpio.OutputPayload{a.paused, a.areas.Value.(AreaType).Color}
	}
	return nil
}

func NewAreaMode() ModeHandler {
	am := AreaMode{
		paused: true,
		areas:  ring.New(0),
	}

	am.RegisterType("Testing", uint32(0x009999))
	am.RegisterType("Testing2", uint32(0x999900))
	return &am
}
