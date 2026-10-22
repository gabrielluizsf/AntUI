package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/svg"
)

// cyberCursorSVG is the pointer's own shape, a plain four-point dart —
// Lucide's own mouse-pointer-2 glyph, which reads as one at a cursor's
// actual size in a way this hand-plotted polygon kept fighting. The fill
// carries the silhouette, since a hollow outline alone all but vanishes at
// sixteen pixels; the pale `currentColor` stroke on top of it is what
// stays this template's own card edge rather than a stock black arrow.
const cyberCursorSVG = `<svg viewBox="0 0 24 24">
	<path d="M4 4 L11.07 21 L13.58 13.61 L21 11.07 Z" fill="#131B33"
		stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/>
</svg>`

// cyberCursorEdge is the pointer's own border colour — [paintPanels]'s pale
// card edge, 0x66EAF4FF, made opaque: a legacy cursor is one colour or none
// at all, nothing in between for a translucent edge to survive as.
const cyberCursorEdge = canvas.Color(0xFFEAF4FF)

// CyberpunkCursor draws the OS mouse pointer this template sets while it is
// active, built to match the card the rest of the screen sits in rather
// than stand apart from it: no solid fill, just the same thin pale edge a
// card is outlined in. win's own [Scale] decides how big it stands, so the
// pointer grows with the window the rest of the template does.
func CyberpunkCursor(win *antui.Window) antui.CursorIcon {
	u := Scale(win)
	img, err := svg.Parse(cyberCursorSVG)
	if err != nil {
		return antui.CursorIcon{}
	}
	size := 24 * u
	cv := img.RenderWith(size, size, cyberCursorEdge)
	if cv == nil {
		return antui.CursorIcon{}
	}
	// The dart's own tip sits at the SVG's (4,4), and a 24-unit viewBox
	// rendered into a 24*u picture is exactly u pixels per unit — no
	// separate scale to carry from the drawing to the hotspot.
	return antui.CursorIcon{Image: cv, HotX: 4 * u, HotY: 4 * u}
}
