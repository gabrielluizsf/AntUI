package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Input is the site's field: a near-black chamfered window with a thin blue
// edge, the same cut the button wears. Hover lights the edge, focus lights a
// full edge and a halo.
func (cyberStyle) Input(p InputParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	cut := buttonRadius(p.H, u)

	body, border := cyberPalette.field, canvas.Fade(v.neon, 51)
	if p.State.Focused {
		border = v.neon
		chamferGlow(cv, halo{x: p.X, y: p.Y, w: p.W, h: p.H, radius: cut, colour: v.neon, strength: int(float64(20) * v.glow), u: u})
	} else if p.State.Hovered {
		border = canvas.Fade(v.neon, 102)
	}
	neonFill(cv, p.X, p.Y, p.W, p.H, cut, u, body, border)
	cyberStyle{}.inputDraw(p)
}

// inputDraw is the rest of [cyberStyle.Input] after its frame, kept apart so
// the frame can take the site's chamfered shape without the text trailing it.
func (cyberStyle) inputDraw(p InputParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	padding := 8 * u
	cursor := min(max(p.Cursor, 0), len(p.Text))
	cursorPx := canvas.TextWidth(p.Text[:cursor]) * u
	scroll := max(cursorPx-(p.W-2*padding), 0)

	savedClip := cv.Clip
	cv.SetClip(canvas.Area{X: p.X + padding, Y: p.Y + u, Width: p.W - 2*padding, Height: p.H - 2*u})

	textY := p.Y + (p.H-textHeight(u))/2
	text := p.Text
	colour := v.text
	if text == "" {
		text, colour = "Seu nome", v.dim
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + padding - scroll), Y: float64(textY)}, text, canvas.TextStyle{Color: colour, Scale: u})
	if p.State.Focused && p.Caret {
		cx := p.X + padding - scroll + cursorPx
		drawCyberCaret(cv, caretGlyph{X: cx, Y: textY, H: textHeight(u), U: u, Color: v.neon})
	}
	cv.Clip = savedClip
}
