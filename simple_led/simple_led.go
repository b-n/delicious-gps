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

var (
	blinkTicker *time.Ticker
	ws          wsEngine
	initialized bool
	lamps       []Lamp
)

func Init(c int) error {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = 255
	opt.Channels[0].LedCount = c

	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		return err
	}
	ws = dev
	lamps = make([]Lamp, c)
	for i := 0; i < c; i++ {
		lamps[i] = Lamp{i, false, false, uint32(0x000000), &ws}
	}

	err = ws.Init()
	if err != nil {
		return err
	}

	initialized = true

	blinkTicker = time.NewTicker(time.Second / 3)
	go func(lamps *[]Lamp) {
		for range blinkTicker.C {
			for i := range *lamps {
				(*lamps)[i].Tick()
			}
		}
	}(&lamps)

	return nil
}

func Color(index int, color uint32) {
	lamps[index].Color(color)
}

func Blink(index int, b bool) {
	lamps[index].Blink(b)
}

func Close() {
	if initialized {
		initialized = false
		blinkTicker.Stop()

		for _, l := range lamps {
			l.Color(uint32(0))
		}
		ws.Wait()
		ws.Fini()
	}
}
