package template

func (cyberStyle) TextAreaTextPos(p TextAreaTextPosParams) int {
	win := p.Win
	u := Scale(win)
	padding := 8 * u
	cursor := min(max(p.Cursor, 0), len(p.Text))
	lines := textLines(p.Text, max((p.W-2*padding)/u, 12))
	return textAreaTextPos(p.Text, lines, padding, p.X, p.Y, p.H, cursor, textHeight(u),
		func(prefix string) int { return textWidth(u, prefix) }, p.MouseX, p.MouseY)
}
