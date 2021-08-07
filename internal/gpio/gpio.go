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

type OutputPayload struct {
	Index int
	Blink bool
	Color uint32
}

var (
	outputChannel chan OutputPayload
	inputChannel  chan simple_button.EventPayload
)

func Init(c chan OutputPayload) chan simple_button.EventPayload {
	outputChannel = c
	inputChannel = make(chan simple_button.EventPayload, 1)
	return inputChannel
}

func OpenOutput(ctx context.Context, done chan bool) error {
	err := simple_led.Init(2)
	if err != nil {
		return err
	}

	simple_led.Color(0, 0xffffff)

	go func() {
		logging.Debug("Opened Output")
		for {
			select {
			case o, ok := <-outputChannel:
				if !ok {
					break
				}
				logging.Debugf("Received output payload %v", o)

				simple_led.Color(o.Index, o.Color)
				simple_led.Blink(o.Index, o.Blink)
			case <-ctx.Done():
				logging.Debug("Stopping Output")

				simple_led.Close()

				done <- true
				return
			}
		}
	}()

	return nil
}

func ListenInput(ctx context.Context, done chan bool) error {
	buttonEvents, err := simple_button.Init()
	if err != nil {
		return err
	}

	initialized := true

	go func() {
		logging.Debug("Watching Input")
		simple_button.RegisterButton(0, 4)
		simple_button.RegisterButton(1, 17)
		simple_button.RegisterButton(2, 27)

		for {
			select {
			case e := <-buttonEvents:
				if !initialized {
					break
				}
				logging.Debugf("Received Input: %v", e)

				// non-blocking send
				select {
				case inputChannel <- e:
				default:
				}
			case <-ctx.Done():
				logging.Debug("Stopping Input")
				initialized = false
				close(inputChannel)
				simple_button.Close()
				done <- true
				return
			}
		}
	}()

	return nil
}
