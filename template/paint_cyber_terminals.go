package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// paintCyberTerminals draws the two end terminals of the trough.
func paintCyberTerminals(cv *canvas.Canvas, v vivid, trackX, cy, trackW, u int) {
	for _, px := range []int{
		trackX,
		trackX + trackW,
	} {
		cv.FillCircle(px, cy, 3*u, cyberPalette.panelAct)
		cv.Circle(px, cy, 3*u, v.edge)
		cv.FillCircle(px, cy, u, v.neon)
	}
}
