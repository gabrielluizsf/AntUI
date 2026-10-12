package template

import "github.com/gabrielluizsf/antui"

// selectKeys walks an open menu with the keyboard: the arrows move the
// highlight, Enter or Space picks the highlighted option, Escape closes. It
// reports whether an option was picked.
func (t *uiTemplate) selectKeys(win *antui.Window, index *int, options []string, open, focused bool) bool {
	if !open || !focused {
		return false
	}
	if win.KeyPressed(antui.KeyUp) && t.selectCursor > 0 {
		t.selectCursor--
	}
	if win.KeyPressed(antui.KeyDown) && t.selectCursor < len(options)-1 {
		t.selectCursor++
	}
	if win.KeyPressed(antui.KeyEnter) || win.KeyPressed(antui.KeySpace) {
		*index = t.selectCursor
		t.selectOpen = 0
		return true
	}
	if win.KeyPressed(antui.KeyEscape) {
		t.selectOpen = 0
	}
	return false
}
