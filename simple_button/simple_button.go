package simple_button

import (
	"os"
	"time"

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
	buttons             []Button
	pollingTicker       *time.Ticker
)

type EventPayload struct {
	Id    uint8
	Pin   uint8
	Event ButtonEvent
}

func notify(b *Button, event ButtonEvent) {
	if !initialized {
		return
	}
	notificationChannel <- EventPayload{
		Id:    b.id,
		Pin:   b.pin,
		Event: event,
	}
}

func Init(pollingInterval time.Duration) (chan EventPayload, error) {
	notificationChannel = make(chan EventPayload)
	buttons = make([]Button, 0)

	// Allow the application to run, even if gpio isn't available (for debugging)
	if _, err := os.Stat("/dev/gpiomem"); os.IsNotExist(err) {
		return notificationChannel, nil
	}

	err := rpio.Open()
	if err != nil {
		return nil, err
	}

	initialized = true

	pollingTicker = time.NewTicker(pollingInterval)

	// TODO: Maybe we don't share memory between goroutines one day
	go func(buttons *[]Button) {
		for range pollingTicker.C {
			for i := range *buttons {
				(*buttons)[i].tick()
			}
		}
	}(&buttons)

	return notificationChannel, nil
}

func RegisterButton(id uint8, gpio_pin uint8) {
	if _, err := os.Stat("/dev/gpiomem"); os.IsNotExist(err) {
		return
	}

	pin := rpio.Pin(gpio_pin)
	pin.Input()
	pin.PullUp()

	buttons = append(buttons, Button{
		pin:             gpio_pin,
		id:              id,
		lastSteadyState: rpio.High,
		currentState:    rpio.High,
		lastState:       rpio.High,
		rpio_pin:        pin,
		timers:          buttonTimers{},
	})
}

func Close() error {
	if initialized {
		initialized = false
		pollingTicker.Stop()
		close(notificationChannel)
		return rpio.Close()
	}
	return nil
}
