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
				t.selectOpen = id
				win.WidgetFocusSet(id)
			} else {
				t.selectOpen = 0
			}
		}
	}
	if open {
		if t.selectMenu(win, x, y, w, h, index, options, focused) {
			open = false
			picked = true
		}
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
