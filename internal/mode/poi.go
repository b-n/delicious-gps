package mode

import (
	"github.com/b-n/delicious-gps/internal/gpio"
	"github.com/b-n/delicious-gps/simple_button"
)

type PoiMode struct{}

func NewPoiMode() ModeHandler {
	pm := PoiMode{}
	return &pm
}

func (p *PoiMode) HandleInput(e simple_button.EventPayload) *gpio.OutputPayload {
	//log POI Datapoint
	return &gpio.OutputPayload{false, uint32(0xff0000)}
}

func (p *PoiMode) HandleLocationEvent(e interface{}) {
	// store in cache
}
