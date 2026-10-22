package template

import "github.com/gabrielluizsf/antui/canvas"

func (builtinStyle) Switch(p SwitchParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	trackW, trackH := switchTrackW*u, switchTrackH*u
	radius := trackH / 2

	// The track itself carries the state — off is a neutral pill, on is the
	// accent — which is what every phone's own switch does; a dot sliding
	// inside a track that never changes colour reads as a dial, not a toggle.
	track := th.SurfaceActive
	switch {
	case p.On && p.State.Hovered:
		track = th.AccentHover
	case p.On:
		track = th.Accent
	case p.State.Hovered:
		track = th.Border
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: trackW, Height: trackH}, radius), track)
	if !p.On {
		// Off is outlined as well as filled, the way a switch nobody has
		// touched yet reads as "nothing chosen" rather than "chosen: no".
		cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: trackW, Height: trackH}, radius), th.Border)
	}
	if p.State.Focused {
		cv.RoundRect(canvas.Round(canvas.Area{X: p.X - focusRing*u, Y: p.Y - focusRing*u, Width: trackW + 2*focusRing*u, Height: trackH + 2*focusRing*u}, radius+focusRing*u), th.Accent)
	}

	// The thumb itself is a second tell for the state, not just the track:
	// off is a small dot adrift in the middle of its half of the track, on
	// is a disc that all but fills it and carries a check.
	thumbR := radius - 6*u
	if p.On {
		thumbR = radius - 2*u
	}
	cx := p.X + radius
	if p.On {
		cx = p.X + trackW - radius
	}
	cy := p.Y + radius
	cv.BoxShadow(canvas.BoxShadowParams{X: cx - thumbR, Y: cy - thumbR, Width: thumbR * 2, Height: thumbR * 2, RX: thumbR, RY: thumbR, DY: u, Blur: 2 * u, Color: th.Shadow})
	cv.FillCircle(canvas.Circle{X: cx, Y: cy, R: thumbR}, th.Surface)
	if p.On {
		drawIconIn(cv, IconCheck, cx-thumbR+2*u, cy-thumbR+2*u, thumbR*2-4*u, th.Accent)
	}

	cv.TextScaled(canvas.Point{X: float64(p.X + trackW + labelGap*u), Y: float64(p.Y + (trackH-textHeight(u))/2)}, p.Label, canvas.TextStyle{Color: th.Text, Scale: u})
}
