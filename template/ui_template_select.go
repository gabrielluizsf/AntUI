package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

// Select draws a dropdown picker. The box shows the currently chosen option;
// clicking it opens or closes a menu of all options right below it. Picking an
// option sets index to its position. It reports when the menu opened, closed,
// or a new option was picked.
func (t *uiTemplate) Select(win *antui.Window, x, y, w, h int, index *int, options []string) event.Event {
	if index == nil || len(options) == 0 {
		return event.Nothing
	}
	if *index < 0 {
		*index = 0
	}
	if *index >= len(options) {
		*index = len(options) - 1
	}

	id := win.WidgetID("template:select", x, y, w, h, "")
	hovered := win.Hovered(x, y, w, h)
	open := t.selectOpen == id
	wasOpen := open
	menuH := h * len(options)

	// A press anywhere outside the box and its menu closes it. This runs
	// first so a click that closes the menu does not also count for the box.
	if open && win.MousePressed(antui.MouseLeft) && !win.Hovered(x, y, w, h+menuH) {
		open = false
	}

	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	picked := false

	// Keyboard, while the menu is open and the box holds the focus: arrows
	// walk the menu, Enter or Space picks it, Escape closes it. It runs before
	// the box's own click activation below, so Enter picks instead of just
	// toggling the menu.
	if open && focused {
		if win.KeyPressed(antui.KeyUp) && t.selectCursor > 0 {
			t.selectCursor--
		}
		if win.KeyPressed(antui.KeyDown) && t.selectCursor < len(options)-1 {
			t.selectCursor++
		}
		if win.KeyPressed(antui.KeyEnter) || win.KeyPressed(antui.KeySpace) {
			*index = t.selectCursor
			open, t.selectOpen = false, 0
			picked = true
		}
		if win.KeyPressed(antui.KeyEscape) {
			open, t.selectOpen = false, 0
		}
	}

	// Open/close with click or Enter/Space. Opening also grabs the keyboard,
	// so the arrows above work from the very first frame the menu is out.
	if wasOpen == open && !picked {
		if win.WidgetKeyActivate(id) || win.WidgetClick(id, hovered) {
			open = !open
			if open {
				t.selectCursor = *index
				t.selectOpen = id
				win.WidgetFocusSet(id)
			} else {
				t.selectOpen = 0
			}
		}
	}

	// Mouse: each menu entry is its own clickable area.
	if open {
		type optState struct {
			x, y  int
			opt   string
			sel   bool
			state State
		}
		opts := make([]optState, 0, len(options))
		for i, opt := range options {
			oy := y + h*(i+1)
			oid := win.WidgetID("template:option", x, oy, w, h, opt)
			oh := win.Hovered(x, oy, w, h)
			sel := i == *index
			if win.WidgetClick(oid, oh) {
				*index = i
				open, t.selectOpen = false, 0
				picked = true
			}
			st := State{Hovered: oh, Focused: focused, On: sel}
			if open && focused && t.selectCursor == i {
				st.Hovered = true
			}
			opts = append(opts, optState{x: x, y: oy, opt: opt, sel: sel, state: st})
		}
		// The menu is drawn dead last, above every widget the flow painted
		// under it, so an open menu never hides behind fields below.
		win.WidgetFrontmost(func() {
			for _, o := range opts {
				t.style.SelectOption(win, o.state, o.x, o.y, w, h, o.opt, o.sel)
			}
		})
	}

	t.style.Select(win, State{
		Hovered: hovered,
		Focused: focused,
		On:      open,
	}, x, y, w, h, options[*index], open)

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
