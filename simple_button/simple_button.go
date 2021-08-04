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

const GPIO_POLLING_INT time.Duration = 200 * time.Microsecond

var (
	initialized         bool
	notificationChannel chan EventPayload
	buttons             []Button
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

func Init() (chan EventPayload, error) {
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

	// TODO: Maybe we don't share memory between goroutines one day
	go func(buttons *[]Button) {
		for {
			if !initialized {
				return
			}
			for i := range *buttons {
				(*buttons)[i].tick()
			}
			time.Sleep(GPIO_POLLING_INT)
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
		close(notificationChannel)
		return rpio.Close()
	}
	return nil
}
