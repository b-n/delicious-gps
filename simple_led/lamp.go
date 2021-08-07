package simple_led

type Lamp struct {
	index    int
	on       bool
	blinking bool
	color    uint32
	ws       *wsEngine
}

func (l *Lamp) render() {
	toColor := uint32(0x000000)
	if l.on {
		toColor = l.color
	}
	(*l.ws).Leds(0)[l.index] = toColor
	(*l.ws).Render()
}

func (l *Lamp) Color(c uint32) {
	l.color = c
	l.render()
}

func (l *Lamp) Blink(b bool) {
	l.blinking = b
	if !l.blinking {
		l.on = true
		l.render()
		return
	}
}

func (l *Lamp) Tick() {
	if !l.blinking {
		return
	}
	l.on = !l.on
	l.render()
}
