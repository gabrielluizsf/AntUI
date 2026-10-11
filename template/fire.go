package template

import "github.com/gabrielluizsf/antui/template/event"

// fire reports an event and, when one is attached, rings it through the
// sound bank before answering.
func (t *uiTemplate) fire(e event.Event) event.Event {
	if e.Ok() && t.sound != nil {
		t.sound.On(e)
	}
	return e
}
