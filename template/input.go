package template

import "github.com/gabrielluizsf/antui/canvas"

func (builtinStyle) Input(p InputParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	padding := 8 * u

	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, th.Radius), th.Surface)
	border := th.Border
	switch {
	case p.State.Focused:
		border = th.Accent
	case p.State.Hovered:
		border = th.TextMuted
	}
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, th.Radius), border)

	cursor := min(max(p.Cursor, 0), len(p.Text))
	cursorPx := canvas.TextWidth(p.Text[:cursor]) * u
	scroll := max(cursorPx-(p.W-2*padding), 0)

	savedClip := cv.Clip
	cv.SetClip(canvas.Area{X: p.X + padding, Y: p.Y + u, Width: p.W - 2*padding, Height: p.H - 2*u})

	textY := p.Y + (p.H-textHeight(u))/2
	text := p.Text
	colour := th.Text
	if text == "" {
		text, colour = "Seu nome", th.TextMuted
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + padding - scroll), Y: float64(textY)}, text, canvas.TextStyle{Color: colour, Scale: u})

	if p.State.Focused && p.Caret {
		cv.FillRect(canvas.Area{X: p.X + padding - scroll + cursorPx, Y: textY - u, Width: u, Height: textHeight(u) + 2*u}, th.Text)
	}
	cv.Clip = savedClip
}
