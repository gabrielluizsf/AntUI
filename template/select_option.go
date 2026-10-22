package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) SelectOption(p SelectOptionParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()

	u := Scale(win)
	// A row in an open menu reads as a list, not a grid of its own buttons:
	// hover gets the same quiet tint a row anywhere else in the app would,
	// and the accent is spent on the bar that marks which one is picked —
	// not on flooding the row behind the label.
	fill := th.Surface
	switch {
	case p.Selected:
		fill = canvas.Fade(th.Accent, 26)
	case p.State.Hovered:
		fill = th.SurfaceHover
	}
	// Each row casts the same shadow; every row but the bottom one has its
	// own fill painted over the part that would fall across the row below,
	// which is what turns a shadow per row into one shadow around the menu.
	cv.BoxShadow(canvas.BoxShadowParams{X: p.X, Y: p.Y, Width: p.W, Height: p.H, DY: 2 * u, Blur: 6 * u, Color: th.Shadow})
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, 0), fill)
	if p.Selected {
		cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: 3 * u, Height: p.H}, th.Accent)
	}

	textY := p.Y + (p.H-textHeight(u))/2
	colour := th.Text
	if p.Selected {
		colour = th.Accent
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + 10*u), Y: float64(textY)}, p.Label, canvas.TextStyle{Color: colour, Scale: u})
}
