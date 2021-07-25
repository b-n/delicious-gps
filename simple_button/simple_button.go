package simple_button

import (
	"context"
	"os"

	"github.com/stianeikeland/go-rpio/v4"
)

const (
	PullUp   = iota
	PullDown = iota
	PullOff  = iota
)

type Button struct {
	pin      uint8
	rpio_pin rpio.Pin
}

func (b *Button) Listen(state chan bool) {
	if *b == (Button{}) {
		return
	}
	go func() {
		for {
			if b.rpio_pin.EdgeDetected() {
				state <- true
			}
		}
	}()
}

var (
	initialized bool
)

func NewSimpleButton(ctx context.Context, gpio_pin uint8) (*Button, error) {
	// Allow the application to run, even if gpio isn't available (for debugging)
	if _, err := os.Stat("/dev/gpiomem"); os.IsNotExist(err) {
		return &Button{}, nil
	}

	if !initialized {
		err := rpio.Open()
		if err != nil {
			return nil, err
		}
		initialized = true
	}

	go func() {
		defer rpio.Close()
		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}()

	pin := rpio.Pin(gpio_pin)
	pin.Input()
	pin.Detect(rpio.RiseEdge)

	butt := Button{
		pin:      gpio_pin,
		rpio_pin: pin,
	}

	return &butt, nil
}
