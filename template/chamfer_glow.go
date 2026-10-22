package template

import "github.com/gabrielluizsf/antui/canvas"

// chamferGlow draws translucent shells around a chamfered rectangle, the
// strongest closest to it, the same way [glow] does for a rounded one — so a
// button's own halo follows the corners it actually cuts instead of filling
// them back in with a square one.
func chamferGlow(cv *canvas.Canvas, h halo) {
	for i := 2; i >= 0; i-- {
		alpha := h.strength / (3 - i)
		if alpha <= 0 {
			continue
		}
		g := (i + 1) * h.u
		pts := chamfer(h.x-g, h.y-g, h.w+g*2, h.h+g*2, h.radius+g)
		fillClip(cv, pts, canvas.Fade(h.colour, alpha))
	}
}
