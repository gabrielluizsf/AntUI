package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) TextArea(p TextAreaParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(p.H, u)

	// Same field as Input: near-black glass, hairline accent frame, focus
	// lights it, hover only lifts the edge a step.
	body, border := cyberPalette.field, canvas.Fade(v.neon, 51)
	if p.State.Focused {
		border = v.neon
		chamferGlow(cv, halo{x: p.X, y: p.Y, w: p.W, h: p.H, radius: radius, colour: v.neon, strength: int(float64(20) * v.glow), u: u})
	} else if p.State.Hovered {
		border = canvas.Fade(v.neon, 102)
	}
	neonFill(cv, p.X, p.Y, p.W, p.H, radius, u, body, border)

	padding := 8 * u
	innerW := p.W - 2*padding
	textW := max(innerW/u, 12)
	cursor := min(max(p.Cursor, 0), len(p.Text))
	placeholder := p.Text == ""
	display := p.Text
	if placeholder {
		display = "Escreva uma anotação..."
	}
	lines := textLines(display, textW)

	lineH := textHeight(u)
	visible := max((p.H-2*padding)/lineH, 1)
	viewTop := textAreaViewTop(cursor, lines, visible)

	savedClip := cv.Clip
	cv.SetClip(canvas.Area{X: p.X + padding, Y: p.Y + padding, Width: innerW, Height: p.H - 2*padding})

	for lif, l := range lines {
		if lif < viewTop || lif >= viewTop+visible {
			continue
		}
		cy := p.Y + padding + (lif-viewTop)*lineH
		line := display[l[0]:l[1]]
		colour := v.text
		if placeholder {
			colour = v.dim
		} else if line == "" {
			colour = v.dim
		}
		cv.TextScaled(canvas.Point{X: float64(p.X + padding), Y: float64(cy)}, line, canvas.TextStyle{Color: colour, Scale: u})
		if p.State.Focused && p.Caret && cursor >= l[0] && cursor <= l[1] {
			off := cursor - l[0]
			cx := p.X + padding + textWidth(u, line[:off])
			drawCyberCaret(cv, caretGlyph{X: cx, Y: cy, H: lineH, U: u, Color: v.neon})
		}
	}
	cv.Clip = savedClip
}
