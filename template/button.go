package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) Button(p ButtonParams) {
	cv := p.Win.Canvas()
	th := p.Win.Theme()
	u := Scale(p.Win)
	radius := min(th.Radius*u, p.H/2)
	fill := th.Accent
	switch {
	case p.State.Pressed:
		fill = canvas.Shade(th.Accent, -0.18)
	case p.State.Hovered:
		fill = th.AccentHover
	}

	// A button sits above the surface it is on; pressing it settles the
	// button onto that surface instead of casting a shadow past it, which is
	// what grounds the motion rather than just darkening the fill.
	if !p.State.Pressed {
		lift := 1
		if p.State.Hovered {
			lift = 2
		}
		cv.BoxShadow(canvas.BoxShadowParams{
			X: p.X, Y: p.Y, Width: p.W, Height: p.H, RX: radius, RY: radius,
			DY: lift * u, Blur: (lift + 2) * u, Color: th.Shadow,
		})
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, radius), fill)

	if p.State.Focused && !p.State.Pressed {
		cv.RoundRect(canvas.Round(canvas.Area{X: p.X - focusRing*u, Y: p.Y - focusRing*u, Width: p.W + 2*focusRing*u, Height: p.H + 2*focusRing*u}, radius+focusRing*u), th.Accent)
	}

	textY := p.Y + (p.H-textHeight(u))/2
	if p.State.Pressed {
		textY += u
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + (p.W-textWidth(u, p.Label))/2), Y: float64(textY)}, p.Label, canvas.TextStyle{Color: th.TextOnAccent, Scale: u})
}
