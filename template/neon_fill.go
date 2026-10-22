package template

import "github.com/gabrielluizsf/antui/canvas"

// neonFill paints a control's body and its border: a dark panel with a thin
// coloured edge, chamfered off the top-right and bottom-left the same way
// the button is — the site's one shape language, not a rounded field next to
// an angular button. What happens under the pointer — who glows, who
// brightens — is the caller's call, because a button and a field answer
// differently.
func neonFill(cv *canvas.Canvas, x, y, w, h, radius, u int, body, border canvas.Color) {
	pts := chamfer(x, y, w, h, radius)
	fillClip(cv, pts, body)
	strokeClip(cv, pts, u, border)
}
