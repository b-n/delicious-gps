package simple_button

import (
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

const DEBOUNCE_MS int64 = 30000

type ButtonEvent uint8

const (
	PRESSED ButtonEvent = iota
	RELEASED
)

var (
	initialized         bool
	notificationChannel chan EventPayload
	buttons             []Button
)

type EventPayload struct {
	Pin   uint8
	Event ButtonEvent
}

type Button struct {
	pin             uint8
	lastSteadyState rpio.State
	lastState       rpio.State
	currentState    rpio.State
	lastDebounceMs  int64
	polling         bool
	rpio_pin        rpio.Pin
}

func (b *Button) Watch(events chan EventPayload) {
	for {
		if !b.polling {
			break
		}

		// Debouncing logic
		b.currentState = b.rpio_pin.Read()

		if b.currentState != b.lastState {
			b.lastDebounceMs = time.Now().UnixNano()
			b.lastState = b.currentState
		}

		if time.Now().UnixNano()-b.lastDebounceMs > DEBOUNCE_MS {
			if b.lastSteadyState == rpio.High && b.currentState == rpio.Low {
				events <- EventPayload{b.pin, PRESSED}
			}
			if b.lastSteadyState == rpio.Low && b.currentState == rpio.High {
				events <- EventPayload{b.pin, RELEASED}
			}

			b.lastSteadyState = b.currentState
		}

	}
}

func (b *Button) Stop() {
	b.polling = false
}

func Init() (chan EventPayload, error) {
	notificationChannel = make(chan EventPayload, 1)

	// Allow the application to run, even if gpio isn't available (for debugging)
	if _, err := os.Stat("/dev/gpiomem"); os.IsNotExist(err) {
		return notificationChannel, nil
	}

	err := rpio.Open()
	if err != nil {
		return nil, err
	}

	initialized = true
	return notificationChannel, nil
}

func Close() error {
	close(notificationChannel)
	if initialized {
		initialized = false
		return rpio.Close()
	}
	return nil
}

func NewSimpleButton(gpio_pin uint8) *Button {
	if _, err := os.Stat("/dev/gpiomem"); os.IsNotExist(err) {
		return &Button{}
	}

	pin := rpio.Pin(gpio_pin)
	pin.Input()
	pin.PullUp()

	butt := Button{
		pin:             gpio_pin,
		lastSteadyState: rpio.High,
		currentState:    rpio.High,
		lastState:       rpio.High,
		polling:         true,
		rpio_pin:        pin,
	}

	go butt.Watch(notificationChannel)

	return &butt
}
