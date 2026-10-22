package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

// Select draws a dropdown picker. The box shows the currently chosen option;
// clicking it opens or closes a menu of all options right below it. Picking an
// option sets index to its position. It reports when the menu opened, closed,
// or a new option was picked.
func (t *uiTemplate) Select(win *antui.Window, area canvas.Area, params SelectParams) event.Event {
	x, y, w, h := area.X, area.Y, area.Width, area.Height
	index := params.Index
	options := params.Options
	if index == nil || len(options) == 0 {
		return event.Nothing
	}
	clampSelectIndex(index, options)

	id := win.WidgetID("template:select", canvas.Area{X: x, Y: y, Width: w, Height: h}, "")
	hovered := win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: h})
	open := t.selectOpen == id
	wasOpen := open

	// A press anywhere outside the box and its menu closes it.
	if open && win.MousePressed(antui.MouseLeft) && !win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: h + h*len(options)}) {
		open = false
	}

	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	picked := t.selectKeys(win, index, options, open, focused)
	if picked {
		open = false
	}

	// Open/close with click or Enter/Space. Opening also grabs the keyboard,
	// so the arrows work from the very first frame the menu is out.
	if wasOpen == open && !picked {
		if win.WidgetKeyActivate(id) || win.WidgetClick(id, hovered) {
			open = !open
			if open {
				t.selectCursor = *index
				win.WidgetFocusSet(id)
			}
		}
	}
	// open may also have closed above — a press outside, Escape, a pick —
	// paths that only ever touch the local variable; selectOpen is the one
	// that has to agree with it, since it alone survives to the next frame,
	// or a menu a click outside just closed springs back open on it.
	if open {
		t.selectOpen = id
	} else if t.selectOpen == id {
		t.selectOpen = 0
	}
	if open {
		// The claim is released before the menu's own rows get a turn at the
		// pointer: each row is its own widget, with its own id the trigger's
		// claim does not cover, and answers for itself regardless.
		win.ReleaseModal(id)
		if t.selectMenu(win, x, y, w, h, index, options, focused) {
			open = false
			picked = true
		}
	}
	// Claimed every frame the menu is open, released the frame it is not —
	// which is what keeps a widget drawn earlier in flow order, underneath
	// where the menu lands, from taking a click meant for whatever the menu
	// is covering it with.
	if open {
		win.ClaimModal(id)
	} else {
		win.ReleaseModal(id)
	}

	t.style.Select(SelectParams{
		Win:   win,
		State: State{Hovered: hovered, Focused: focused, On: open},
		X:     x, Y: y, W: w, H: h,
		Value: options[*index],
		Open:  open,
	})

	return t.selectResult(wasOpen, open, picked, index)
}
