package gpio

import (
	"context"

	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/b-n/delicious-gps/simple_button"
	"github.com/b-n/delicious-gps/simple_led"
)

var (
	button *simple_button.Button
)

func OpenOutput(ctx context.Context, done chan bool, startState uint8) (chan uint8, error) {
	outputChannel := make(chan uint8)
	state := startState

	led, err := simple_led.NewSimpleLED()
	if err != nil {
		return nil, err
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

	return outputChannel, nil
}

func ListenInput(ctx context.Context, done chan bool, inputChannel chan uint8) error {
	button, err := simple_button.NewSimpleButton(4)
	if err != nil {
		return err
	}

	go func() {
		logging.Debug("Watching Input")
		buttonReleased := make(chan bool)
		button.Listen(buttonReleased)
		for {
			select {
			case e := <-buttonEvents:
				logging.Debugf("Received Input: %v", e)

				// TODO: buffer the events to main
				select {
				case inputEvents <- e.event:
				default:
				}
			case <-ctx.Done():
				logging.Debug("Stopping Input")
				button.Close()

				done <- true
				return
			}
		}
	}()

	return nil
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
