package template

import "github.com/gabrielluizsf/antui"

func (cyberStyle) Button(win *antui.Window, s State, x, y, w, h int, label string) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(h, u)

	// Idle is a quiet panel with a blue edge. Attention lights it: the border
	// turns neon and a soft halo grows around it, larger while it is held.
	body, border := cyberPalette.panel, v.edge
	reach := 0
	switch {
	case s.Pressed:
		body, border, reach = cyberPalette.panelAct, v.neon, 30
	case s.Hovered:
		body, border, reach = cyberPalette.panelHov, v.neon, 18
	case s.Focused:
		border, reach = v.neon, 16
	}
	if reach > 0 {
		glow(cv, halo{x: x, y: y, w: w, h: h, radius: radius, colour: v.neon, strength: int(float64(reach) * v.glow), u: u})
	}
	neonFill(cv, x, y, w, h, radius, body, border)
	cv.TextScaled(x+(w-textWidth(u, label))/2, y+(h-textHeight(u))/2, label, v.text, u)
}
