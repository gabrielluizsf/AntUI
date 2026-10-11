package template

// interactionState reads the state the shared interaction machinery assigns a
// widget: hovered when the pointer is in its hover box, focused when the
// keyboard arrived, pressed when a press is landing on it, and active while
// it is the one being held. The identity box matches the widget's own
// WidgetID — for focus to agree — and idW/idH are those calling arguments;
// hovW/hovH are the box the widget itself counts a hover in. extra carries
// the states the widget decided itself, the On of a checked box included.
func (c *CSS) interactionState(kind string, x, y, idW, idH int, label string, hovW, hovH int, extra State) State {
	s := extra
	id := c.win.WidgetID("template:"+kind, x, y, idW, idH, label)
	hovered := c.win.Hovered(x, y, hovW, hovH)
	s.Hovered = hovered
	s.Focused = c.win.WidgetFocus() == id
	s.Active = c.win.WidgetActive() == id
	if s.Active && hovered {
		s.Pressed = true
	}
	return s
}
