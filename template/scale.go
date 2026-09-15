package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// ScaleBase is the reference window edge that Scale draws for at scale 1:
// a 360x360 window or smaller. Everything a template draws — boxes, padding,
// icons, borders — is a coordinate unit times the scale, so a screen designed
// and tested at 360 wide scales up to 4K and stays as airy as it was at the
// reference size.
const ScaleBase = 360

// clampInt keeps v inside [lo, hi].
func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Scale is how many drawing units one reference-scale unit is at the given
// window. It follows the smaller edge, so a phone in portrait and in
// landscape agree on the size of a button. The floor of 1 keeps a tiny
// window from shrinking a control so far it cannot be pressed; the cap of 8
// keeps an enormous display from drawing grid lines a mile apart.
func Scale(win *antui.Window) int {
	if win == nil {
		return 1
	}
	smaller := min(win.Width(), win.Height())
	return clampInt(smaller/ScaleBase, 1, 8)
}

// ScaleOfContext is [Scale] for a drawing Context that has a window.
func (c *Context) Scale() int { return Scale(c.win) }

// textWidth measures a string at the scale a Context draws text at — the
// built-in font times its scale. Screens use it where nothing shrinks to fit:
// to centre a button whose width follows its label, for instance.
func textWidth(u int, text string) int { return canvas.TextWidth(text) * u }

// textHeight is the line height of scaled text, drawn in the canvas's
// current default face — the built-in face at 16 px per line, a system face
// at whatever height that font actually uses.
func textHeight(u int) int { return canvas.TextHeight() * u }

// TextWidth is textWidth on the Context's own scale.
func (c *Context) TextWidth(text string) int { return textWidth(c.Scale(), text) }