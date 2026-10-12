package template

import "github.com/gabrielluizsf/antui/template/event"

// selectResult reports what the frame did: a new option picked, or the menu
// opening or closing.
func (t *uiTemplate) selectResult(wasOpen, open, picked bool, index *int) event.Event {
	if picked {
		return t.fire(event.Event{Component: event.Select, Kind: event.Pick, Step: *index})
	}
	if wasOpen != open {
		kind := event.Close
		if open {
			kind = event.Open
		}
		return t.fire(event.Event{Component: event.Select, Kind: kind})
	}
	return event.Nothing
}
