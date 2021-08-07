package mode

import (
	"time"

	"github.com/b-n/delicious-gps/internal/gpio"
	"github.com/b-n/delicious-gps/simple_button"
)

type PoiMode struct {
	lastLocationEvent interface{}
	lastEvent         time.Time
}

func NewPoiMode() ModeHandler {
	pm := PoiMode{}
	return &pm
}

func (p *PoiMode) HandleInput(e simple_button.EventPayload) {
	if e.Event != simple_button.CLICK {
		return
	}
	if p.lastEvent.Before(time.Now().Add(-5 * time.Second)) {
		writeDisplayChannel(gpio.OutputPayload{0, false, uint32(0xff0000)})
	}

	writeDataChannel(p.lastLocationEvent, POI, 0)
	writeDisplayChannel(gpio.OutputPayload{0, false, uint32(0xffffff)})
}

func (p *PoiMode) HandleLocationEvent(e interface{}) {
	p.lastLocationEvent = e
	p.lastEvent = time.Now()
}

func (a *PoiMode) Activate() {
	writeDisplayChannel(gpio.OutputPayload{0, true, uint32(0xffffff)})
}
