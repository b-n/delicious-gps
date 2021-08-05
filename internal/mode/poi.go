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

func (p *PoiMode) HandleInput(e simple_button.EventPayload) *gpio.OutputPayload {
	if p.lastEvent.Before(time.Now().Add(-5 * time.Second)) {
		return &gpio.OutputPayload{false, uint32(0xff0000)}
	}

	writeDataChannel(p.lastLocationEvent, POI, 0)
	return &gpio.OutputPayload{false, uint32(0xffffff)}
}

func (p *PoiMode) HandleLocationEvent(e interface{}) {
	p.lastLocationEvent = e
	p.lastEvent = time.Now()
}
