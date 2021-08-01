package simple_button

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

const GPIO_POLLING_INT time.Duration = 20 * time.Nanosecond

const DEBOUNCE_MS int64 = 30 * int64(time.Millisecond)
const CLICK_MS int64 = 10 * int64(time.Millisecond)
const DBL_CLICK_DUR int64 = 300 * int64(time.Millisecond)

type buttonTimers struct {
	debounce   int64
	pressed    int64
	click      int64
	clickCount int8
}

type Button struct {
	pin             uint8
	lastSteadyState rpio.State
	lastState       rpio.State
	currentState    rpio.State
	rpio_pin        rpio.Pin
	timers          buttonTimers
}

func (b *Button) Listen() {
	for {
		if !initialized {
			break
		}

		b.currentState = b.rpio_pin.Read()

		//track last known state and set debounce timer
		now := time.Now().UnixNano()
		if b.currentState != b.lastState {
			b.timers.debounce = now
			b.lastState = b.currentState
		}

		// get real state (after debounce)
		if now-b.timers.debounce > DEBOUNCE_MS {
			if b.lastSteadyState == rpio.High && b.currentState == rpio.Low {
				b.timers.pressed = now
				notify(b.pin, ON)
			}
			if b.lastSteadyState == rpio.Low && b.currentState == rpio.High {
				notify(b.pin, OFF)
			}

			b.lastSteadyState = b.currentState
		}

		// Clickevent logic
		// - off-on needs to be long enough to register a click (increments count if true)
		// - will wait for multi click events, otherwise fires event (based on count)
		if b.timers.pressed > 0 && b.currentState == rpio.High && now-b.timers.pressed > CLICK_MS {
			fmt.Printf("parsing click %d\n", b.timers.clickCount)
			b.timers.pressed = 0
			b.timers.click = now
			b.timers.clickCount = b.timers.clickCount + 1
		}

		// wait for multi click events
		if b.timers.click > 0 && now-DBL_CLICK_DUR > b.timers.click {
			fmt.Printf("sending event %d\n", b.timers.clickCount)
			if b.timers.clickCount > 1 {
				notify(b.pin, DBL_CLICK)
			} else {
				notify(b.pin, CLICK)
			}

			b.timers.click = 0
			b.timers.clickCount = 0
		}

		time.Sleep(GPIO_POLLING_INT)
	}
}
