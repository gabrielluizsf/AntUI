package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Input(p InputParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(p.H, u)

	body, border := cyberPalette.panel, v.edge
	if p.State.Focused {
		border = v.neon
		glow(cv, halo{x: p.X, y: p.Y, w: p.W, h: p.H, radius: radius, colour: v.neon, strength: int(float64(20) * v.glow), u: u})
	} else if p.State.Hovered {
		body, border = cyberPalette.panelHov, v.neon
	}
	neonFill(cv, p.X, p.Y, p.W, p.H, radius, body, border)

	padding := 8 * u
	cursor := min(max(p.Cursor, 0), len(p.Text))
	cursorPx := canvas.TextWidth(p.Text[:cursor]) * u
	scroll := max(cursorPx-(p.W-2*padding), 0)

	savedClip := cv.Clip
	cv.SetClip(canvas.Area{X: p.X + padding, Y: p.Y + u, Width: p.W - 2*padding, Height: p.H - 2*u})

	textY := p.Y + (p.H-textHeight(u))/2
	colour := v.text
	if p.Text == "" {
		colour = v.dim
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + padding - scroll), Y: float64(textY)}, p.Text, canvas.TextStyle{Color: colour, Scale: u})
	if p.State.Focused && p.Caret {
		cx := p.X + padding - scroll + cursorPx
		cv.FillRect(canvas.Area{X: cx, Y: textY, Width: 2 * u, Height: textHeight(u)}, v.neon)
		cv.FillRect(canvas.Area{X: cx, Y: textY + textHeight(u), Width: 2 * u, Height: 2 * u}, canvas.Fade(v.neon, 150))
	}
	cv.Clip = savedClip
}
