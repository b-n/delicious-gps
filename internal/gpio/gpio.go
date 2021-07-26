package gpio

import (
	"context"

	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/b-n/delicious-gps/simple_button"
	"github.com/b-n/delicious-gps/simple_led"
)

type InputEvent uint8

const (
	PRESSED InputEvent = iota
	RELEASED
)

var (
	button *simple_button.Button
)

func OpenOutput(ctx context.Context, done chan bool, c chan uint8, startState uint8) error {
	outputChannel := c
	state := startState

	led, err := simple_led.NewSimpleLED()
	if err != nil {
		return nil
	}

	led.Color(uint32(0x0000ff))

	go func() {
		logging.Debug("Opened Output")
		for {
			select {
			case s := <-outputChannel:
				if s != state {
					state = s
					led.Color(colorFromState(state))
				}
			case <-ctx.Done():
				logging.Debug("Stopping Output")

				led.Color(uint32(0x000000))
				led.Close()

				done <- true
				return
			}
		}
	}()

	return nil
}

func ListenInput(ctx context.Context, done chan bool) (chan InputEvent, error) {
	inputEvents := make(chan InputEvent, 1)

	buttonEvents, err := simple_button.Init()
	if err != nil {
		return nil, err
	}

	initialized := true

	go func() {
		logging.Debug("Watching Input")
		simple_button.NewSimpleButton(4)

		for {
			select {
			case e := <-buttonEvents:
				if !initialized {
					break
				}
				logging.Debugf("Received Input: %v", e)
				if e.Event == simple_button.PRESSED {
					// TODO: buffer the events to main
					select {
					case inputEvents <- PRESSED:
					default:
					}
				}
			case <-ctx.Done():
				logging.Debug("Stopping Input")
				initialized = false
				close(inputEvents)
				simple_button.Close()
				done <- true
				return
			}
		}
	}()

	return inputEvents, nil
}

func colorFromState(value uint8) uint32 {
	// Remove MSB
	switch value & 127 {
	case 0:
		return uint32(0x0000ff)
	case 1:
		return uint32(0xff9900)
	case 2:
		return uint32(0xffff00)
	case 3:
		return uint32(0x00ff00)
	case 4:
		return uint32(0x00ff66)
	default:
		return uint32(0xff0000)
	}
}
