package template

import "github.com/gabrielluizsf/antui/canvas"

func (cyberStyle) Button(p ButtonParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(p.H, u)

	// Idle is a quiet panel with a blue edge. Attention lights it: the border
	// turns neon and a soft halo grows around it, larger while it is held.
	body, border := cyberPalette.panel, v.edge
	reach := 0
	switch {
	case p.State.Pressed:
		body, border, reach = cyberPalette.panelAct, v.neon, 30
	case p.State.Hovered:
		body, border, reach = cyberPalette.panelHov, v.neon, 18
	case p.State.Focused:
		border, reach = v.neon, 16
	}
	if reach > 0 {
		glow(cv, halo{x: p.X, y: p.Y, w: p.W, h: p.H, radius: radius, colour: v.neon, strength: int(float64(reach) * v.glow), u: u})
	}
	neonFill(cv, p.X, p.Y, p.W, p.H, radius, body, border)
	cv.TextScaled(canvas.Point{X: float64(p.X + (p.W-textWidth(u, p.Label))/2), Y: float64(p.Y + (p.H-textHeight(u))/2)}, p.Label, canvas.TextStyle{Color: v.text, Scale: u})
}
