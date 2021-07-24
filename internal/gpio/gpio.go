package gpio

import (
	"context"
	"fmt"

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
	if button, err = simple_button.NewSimpleButton(ctx, 18); err != nil {
		return nil, err
	}
	go watchInputChannel(ctx, inputChannel)

	return outputChannel, nil
}

func watchInputChannel(ctx context.Context, inputState chan uint8) {
	buttonReleased := make(chan bool)
	button.Listen(buttonReleased)
	for {
		select {
		case <-buttonReleased:
			fmt.Printf("Button Press Received")
			inputState <- 1
		case <-ctx.Done():
			return
		}
	}
}

func watchOutputChannel(ctx context.Context, newState chan uint8) {
	defer led.Close()
	for {
		select {
		case s := <-newState:
			if s != state {
				fmt.Printf("New state %d received", s)
				state = s
				led.Color(colorFromState(state))
			}
		case <-ctx.Done():
			return
		}
	}
}

func colorFromState(value uint8) uint32 {
	// Remove MSB
	switch value << 1 >> 1 {
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
