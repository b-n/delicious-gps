package gpio

import (
	"context"

	"github.com/b-n/delicious-gps/internal/logging"
	"github.com/b-n/delicious-gps/simple_button"
	"github.com/b-n/delicious-gps/simple_led"
)

var (
	state  uint8
	led    *simple_led.Lamp
	button *simple_button.Button
)

func Open(ctx context.Context, inputChannel chan uint8) (chan uint8, error) {
	outputChannel := make(chan uint8)
	state = 255
	err := error(nil)

	// Setup output
	if led, err = simple_led.NewSimpleLED(); err != nil {
		return nil, err
	}
	led.Color(uint32(0x0000ff))
	go watchOutputChannel(ctx, outputChannel)

	// Setup input
	if button, err = simple_button.NewSimpleButton(ctx, 4); err != nil {
		return nil, err
	}
	go watchInputChannel(ctx, inputChannel)

	return outputChannel, nil
}

func watchInputChannel(ctx context.Context, inputState chan uint8) {
	logging.Debug("Watching for button events (input)")
	buttonReleased := make(chan bool)
	button.Listen(buttonReleased)
	for {
		select {
		case <-buttonReleased:
			logging.Debug("Button Release received")
			inputState <- 1
		case <-ctx.Done():
			logging.Debug("Stopping button events watching")
			return
		}
	}
}

func watchOutputChannel(ctx context.Context, newState chan uint8) {
	logging.Debug("Watching for status changes (output)")
	defer led.Close()
	for {
		select {
		case s := <-newState:
			if s != state {
				logging.Debugf("New state %d received. Current %d", s, state)
				state = s
				led.Color(colorFromState(state))
			}
		case <-ctx.Done():
			logging.Debug("Stopping status changes watching")
			return
		}
	}
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
