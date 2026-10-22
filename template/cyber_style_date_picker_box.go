package template

import "github.com/gabrielluizsf/antui/canvas"

func (cyberStyle) DatePickerBox(p DatePickerBoxParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(p.H, u)

	// Same field language as Select: open and focus light the frame, hover
	// only lifts the edge.
	body, border, reach := cyberPalette.field, canvas.Fade(v.neon, 51), 0
	if p.State.Focused || p.Open {
		border, reach = v.neon, 20
	} else if p.State.Hovered {
		body, border = cyberPalette.panelHov, canvas.Fade(v.neon, 102)
	}
	if reach > 0 {
		chamferGlow(cv, halo{x: p.X, y: p.Y, w: p.W, h: p.H, radius: radius, colour: v.neon, strength: int(float64(reach) * v.glow), u: u})
	}
	neonFill(cv, p.X, p.Y, p.W, p.H, radius, u, body, border)

	padding := 8 * u
	colour := v.text
	if p.Value == "" {
		colour = v.dim
	}
	savedClip := cv.Clip
	cv.SetClip(canvas.Area{X: p.X + padding, Y: p.Y, Width: p.W - 2*padding, Height: p.H})
	cv.TextScaled(canvas.Point{X: float64(p.X + padding), Y: float64(p.Y + (p.H-textHeight(u))/2)}, p.Value, canvas.TextStyle{Color: colour, Scale: u})
	cv.Clip = savedClip

	// The caret that says "there is a calendar below", flipped the way the
	// dropdown's is while the picker is open.
	name, caretTint := IconCaret, v.dim
	if p.Open {
		name, caretTint = IconCaretUp, v.neon
	}
	caretBox := 10 * u
	drawIcon(cv, name, p.X+p.W-padding-3*u-caretBox/2, p.Y+(p.H-caretBox)/2, caretBox, caretBox, caretTint)
}
