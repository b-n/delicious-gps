package simple_led

import (
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

type wsEngine interface {
	Init() error
	Render() error
	SetBrightness(int, int)
	Fini()
	Wait() error
	Leds(channel int) []uint32
}

type Lamp struct {
	on        bool
	blinking  bool
	blinkDone chan bool
	ws        wsEngine
}

func (l *Lamp) Color(c uint32) error {
	l.ws.Leds(0)[0] = c
	return l.ws.Render()
}

func (l *Lamp) Blink(b bool) {
	if l.blinking == b {
		return
	}

	l.blinking = b

	if !b {
		// Send non-blocking done signal
		select {
		case l.blinkDone <- true:
		default:
		}
		return
	}

	go func() {
		ticker := time.NewTicker(time.Second / 3)
		defer ticker.Stop()
		for {
			select {
			case <-l.blinkDone:
				l.ws.SetBrightness(0, 255)
				l.ws.Render()
				l.ws.Wait()
				return
			case <-ticker.C:
				l.on = !l.on
				var brightness int
				if l.on {
					brightness = 255
				} else {
					brightness = 0
				}
				l.ws.SetBrightness(0, brightness)
				l.ws.Render()
			}
		}
	}()
}

func (l *Lamp) Close() {
	if l.blinking {
		l.blinkDone <- true
	}
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

	l := Lamp{true, false, make(chan bool), dev}
	err = l.ws.Init()
	if err != nil {
		return nil, err
	}

	return &l, nil
}
