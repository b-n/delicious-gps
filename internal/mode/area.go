package mode

import (
	"container/ring"

	"github.com/b-n/delicious-gps/internal/gpio"
	"github.com/b-n/delicious-gps/simple_button"
)

type AreaType struct {
	Id    int
	Color uint32
}

type AreaMode struct {
	paused bool
	areas  *ring.Ring
}

func (a *AreaMode) RegisterType(color uint32) {
	var id = 0
	if a.areas != nil {
		id = a.areas.Len()
	}
	j := ring.New(1)
	j.Value = AreaType{Id: id, Color: color}

	if a.areas == nil {
		a.areas = j
	} else {
		a.areas = a.areas.Link(j)
	}
}

func (a *AreaMode) HandleLocationEvent(e interface{}) {
	if !a.paused {
		writeDataChannel(e, AREA, a.areas.Value.(AreaType).Id)
	}
}

func (a *AreaMode) HandleInput(e simple_button.EventPayload) {
	switch e.Event {
	case simple_button.DBL_CLICK:
		a.paused = !a.paused
		writeDisplayChannel(gpio.OutputPayload{a.paused, a.areas.Value.(AreaType).Color})
	case simple_button.CLICK:
		a.areas = a.areas.Next()
		writeDisplayChannel(gpio.OutputPayload{a.paused, a.areas.Value.(AreaType).Color})
	}
}

func (a *AreaMode) Activate() {
	writeDisplayChannel(gpio.OutputPayload{a.paused, a.areas.Value.(AreaType).Color})
}

func NewAreaMode() ModeHandler {
	am := AreaMode{
		paused: true,
		areas:  ring.New(0),
	}

	am.RegisterType(uint32(0x990000))
	am.RegisterType(uint32(0x009900))
	am.RegisterType(uint32(0x009999))
	am.RegisterType(uint32(0x999900))
	return &am
}
