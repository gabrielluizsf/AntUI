package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// Switch draws a flip toggle like a checkbox but a sliding thumb. It flips
// value when clicked and reports whether it flipped this frame.
func (c *CSS) Switch(label string, value *bool) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleSwitch, label)
	if !ok {
		return event.Nothing
	}
	u := c.style.u()
	trackW, trackH := switchTrackW*u, switchTrackH*u
	cy := y + (h-trackH)/2
	on := value != nil && *value
	s := c.interactionState(widgetBox{"switch", label, x, cy, trackW, trackH, trackW + labelGap*u + textWidth(u, label), trackH}, State{On: on})
	st = c.style.entryStyle(c.style.cur, css.RoleSwitch, s)
	if !interactive(st) {
		c.paint(st, css.RoleSwitch, label, x, y, w, h, func() {
			c.style.Switch(SwitchParams{
				Win:   c.win,
				State: s,
				X:     x, Y: cy,
				Label: label,
				On:    on,
			})
		})
		return event.Nothing
	}
	var e event.Event
	c.paint(st, css.RoleSwitch, label, x, y, w, h, func() {
		e = c.ui.Switch(c.win, canvas.Area{X: x, Y: cy}, antui.CheckboxParams{Label: label, Value: value})
	})
	return e
}
