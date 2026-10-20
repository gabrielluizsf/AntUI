package template

import "github.com/gabrielluizsf/antui/canvas"

func (cyberStyle) Button(p ButtonParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	// The site's primary button: a solid accent slab with dark ink, a chamfer
	// off the top-right and bottom-left corners, and no border. Hover bleaches
	// it white inside a soft halo, press does the same louder.
	body := v.neon
	reach := 0
	switch {
	case p.State.Pressed:
		body, reach = canvas.White, 30
	case p.State.Hovered:
		body, reach = canvas.White, 20
	case p.State.Focused:
		reach = 16
	}
	if reach > 0 {
		glow(cv, halo{x: p.X, y: p.Y, w: p.W, h: p.H, radius: 0, colour: v.neon, strength: int(float64(reach) * v.glow), u: u})
	}
	pts := chamfer(p.X, p.Y, p.W, p.H, 10*u)
	fillClip(cv, pts, body)
	cv.TextScaled(canvas.Point{X: float64(p.X + (p.W-textWidth(u, p.Label))/2), Y: float64(p.Y + (p.H-textHeight(u))/2)}, p.Label,
		canvas.TextStyle{Color: cyberPalette.bg, Scale: u, Bold: true, LetterSpacing: u})
}

// chamfer is the site's button outline: a rectangle with the top-right and
// bottom-left corners cut by cut pixels.
func chamfer(x, y, w, h, cut int) []canvas.Point {
	cut = min(cut, min(w, h)/2)
	return []canvas.Point{
		{X: float64(x), Y: float64(y)},
		{X: float64(x + w - cut), Y: float64(y)},
		{X: float64(x + w), Y: float64(y + cut)},
		{X: float64(x + w), Y: float64(y + h)},
		{X: float64(x + cut), Y: float64(y + h)},
		{X: float64(x), Y: float64(y + h - cut)},
	}
}
