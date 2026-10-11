package template

import "github.com/gabrielluizsf/antui"

// vividFor reads the window's opacity and works out how loudly to draw. The
// boost is 255/opacity, one at full opacity and three at a third — the
// brightest ink a colour can survive.
func vividFor(win *antui.Window) vivid {
	o := int(win.Opacity())
	f := 1.0
	if o > 0 {
		f = 255.0 / float64(o)
	}
	if f > 3 {
		f = 3
	}
	return vivid{
		text: lift(cyberPalette.text, f),
		neon: lift(cyberPalette.neon, f),
		dim:  lift(cyberPalette.textDim, f),
		edge: lift(cyberPalette.edge, f),
		glow: f,
	}
}
