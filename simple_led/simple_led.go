package simple_led

import ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"

type wsEngine interface {
	Init() error
	Render() error
	Fini()
	Leds(channel int) []uint32
}

type Lamp struct {
	ws wsEngine
}

func (l *Lamp) Color(c uint32) error {
	l.ws.Leds(0)[0] = c
	return l.ws.Render()
}

func (l *Lamp) Close() {
	l.ws.Fini()
}

func NewSimpleLED() (*Lamp, error) {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = 255
	opt.Channels[0].LedCount = 1

	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		return nil, err
	}

	l := Lamp{ws: dev}
	err = l.ws.Init()
	if err != nil {
		return nil, err
	}

	return &l, nil
}
