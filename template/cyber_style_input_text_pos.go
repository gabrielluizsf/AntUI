package template

func (cyberStyle) InputTextPos(p InputTextPosParams) int {
	win := p.Win
	u := Scale(win)
	padding := 8 * u
	cursor := min(max(p.Cursor, 0), len(p.Text))
	cursorPx := textWidth(u, p.Text[:cursor])
	scroll := max(cursorPx-(p.W-2*padding), 0)
	return caretAt(func(prefix string) int { return textWidth(u, prefix) }, p.Text, p.X+padding-scroll, p.MouseX)
}
