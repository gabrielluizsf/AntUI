package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Switch(p SwitchParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	trackW, trackH := switchTrackW*u, switchTrackH*u
	track := canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: trackW, Height: trackH}, trackH/2)

	if p.State.Focused {
		glow(cv, halo{x: p.X, y: p.Y, w: trackW, h: trackH, radius: trackH, colour: v.neon, strength: int(float64(12) * v.glow), u: u})
	}
	// Off is a dark rail; on fills it with an accent wash. Hover only lifts
	// the edge, checked lights it.
	body, border := cyberPalette.field, v.edge
	if p.State.Hovered {
		border = canvas.Fade(v.neon, 110)
	}
	if p.On {
		body, border = canvas.Fade(v.neon, 40), v.neon
	}
	cv.FillRoundRect(track, body)
	cv.RoundRect(track, border)

	colour := v.dim
	if p.On {
		colour = v.neon
	}
	cx := p.X + trackH/2
	if p.On {
		cx = p.X + trackW - trackH/2
	}
	cv.FillCircle(canvas.Circle{X: cx, Y: p.Y + trackH/2, R: trackH/2 - 3*u}, colour)
	cv.TextScaled(canvas.Point{X: float64(p.X + trackW + labelGap*u), Y: float64(p.Y + (trackH-textHeight(u))/2)}, p.Label, canvas.TextStyle{Color: v.text, Scale: u})
}
