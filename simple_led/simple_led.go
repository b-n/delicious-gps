package simple_led

import (
	"sync"
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

type renderReq struct {
	Index int
	Color uint32
}

var (
	blinkTicker *time.Ticker
	ws          wsEngine
	initialized bool
	lamps       []Lamp
	renderChan  chan renderReq
	done        sync.WaitGroup
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
		lamps[i] = Lamp{i, false, false, uint32(0x000000)}
	}

	err = ws.Init()
	if err != nil {
		return err
	}

	initialized = true

	renderChan = make(chan renderReq, 10)

	go func() {
		done.Add(1)
		for r := range renderChan {
			ws.Leds(0)[r.Index] = r.Color
			ws.Render()
			ws.Wait()
		}
		done.Done()
	}()

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
		close(renderChan)
		done.Wait()
		for _, l := range lamps {
			l.Color(uint32(0))
		}
		ws.Render()
		ws.Wait()
		ws.Fini()
	}
}
