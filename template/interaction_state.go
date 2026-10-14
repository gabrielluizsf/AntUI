package template

import "github.com/gabrielluizsf/antui/canvas"

// interactionState reads the state the shared interaction machinery assigns a
// widget: hovered when the pointer is in its hit box, focused when the keyboard
// arrived, pressed when a press is landing on it, and active while it is the one
// being held. extra carries the states the widget decided itself, the On of a
// checked box included.
func (c *CSS) interactionState(b widgetBox, extra State) State {
	s := extra
	id := c.win.WidgetID("template:"+b.kind, b.x, b.y, b.idW, b.idH, b.label)
	hovered := c.win.Hovered(canvas.Area{X: b.x, Y: b.y, Width: b.hitW, Height: b.hitH})
	s.Hovered = hovered
	s.Focused = c.win.WidgetFocus() == id
	s.Active = c.win.WidgetActive() == id
	if s.Active && hovered {
		s.Pressed = true
	}
	return s
}
