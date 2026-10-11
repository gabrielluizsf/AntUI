package template

import "github.com/gabrielluizsf/antui"

func (cyberStyle) InputTextPos(win *antui.Window, s State, x, y, w, h int, text string, cursor, mx, my int) int {
	u := Scale(win)
	padding := 8 * u
	cursor = min(max(cursor, 0), len(text))
	cursorPx := textWidth(u, text[:cursor])
	scroll := max(cursorPx-(w-2*padding), 0)
	return caretAt(func(prefix string) int { return textWidth(u, prefix) }, text, x+padding-scroll, mx)
}
