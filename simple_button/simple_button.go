package simple_button

import (
	"os"

	"github.com/stianeikeland/go-rpio/v4"
)

type ButtonEvent uint8

const (
	ON ButtonEvent = iota
	OFF
	CLICK
	DBL_CLICK
)

var (
	initialized         bool
	notificationChannel chan EventPayload
)

type EventPayload struct {
	Pin   uint8
	Event ButtonEvent
}

func notify(pin uint8, event ButtonEvent) {
	if !initialized {
		return
	}
	notificationChannel <- EventPayload{
		Pin:   pin,
		Event: event,
	}
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

func Watch(gpio_pin uint8) {
	if _, err := os.Stat("/dev/gpiomem"); os.IsNotExist(err) {
		return
	}

	pin := rpio.Pin(gpio_pin)
	pin.Input()
	pin.PullUp()

	butt := Button{
		pin:             gpio_pin,
		lastSteadyState: rpio.High,
		currentState:    rpio.High,
		lastState:       rpio.High,
		rpio_pin:        pin,
		timers:          buttonTimers{0, 0, 0, 0},
	}

	go butt.Listen()
}
