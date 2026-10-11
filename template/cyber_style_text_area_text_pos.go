package template

import "github.com/gabrielluizsf/antui"

func (cyberStyle) TextAreaTextPos(win *antui.Window, s State, x, y, w, h int, text string, cursor, mx, my int) int {
	u := Scale(win)
	padding := 8 * u
	cursor = min(max(cursor, 0), len(text))
	lines := textLines(text, max((w-2*padding)/u, 12))
	return textAreaTextPos(text, lines, padding, x, y, h, cursor, textHeight(u),
		func(prefix string) int { return textWidth(u, prefix) }, mx, my)
}
